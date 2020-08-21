package executor

import (
	"fmt"
	"io"
	"strconv"
	"sync"

	dataLoaderPkg "github.com/hmalatini/request_executor_script/src/api/data_loader"
	"github.com/hmalatini/request_executor_script/src/api/logger"
	"github.com/hmalatini/request_executor_script/src/api/processor"
	resultWriterPkg "github.com/hmalatini/request_executor_script/src/api/result_writer"
)

const className = "ParallelExecutor"

type ParallelExecutor struct {
	dataLoader   dataLoaderPkg.CsvDataLoader
	processor    processor.RequestProcessor
	resultWriter resultWriterPkg.CsvResultWriter
}

func NewParallelExecutor(dataLoader dataLoaderPkg.CsvDataLoader,
	processor processor.RequestProcessor,
	resultWriter resultWriterPkg.CsvResultWriter) *ParallelExecutor {

	return &ParallelExecutor{
		dataLoader:   dataLoader,
		processor:    processor,
		resultWriter: resultWriter,
	}
}

func (e *ParallelExecutor) Execute(routines int, flushNumber int) {
	// In this channel, there will contain all the data loader by the file input. Whenever it start to receive data,
	// it will process the data with an arbitrary goroutine. IT IS A BUFFER CHANNEL
	dataReceiveCh := make(chan []string, flushNumber)
	// Channel for write the results of the processor, and to consume for write in the output file.
	// IT IS A BUFFER CHANNEL
	resultsCh := make(chan []string, flushNumber)
	logger.LogTrace(className, "Channels were created")

	// Declare the wait group, so the program can know when the goroutines are finished the work
	generalGroup := &sync.WaitGroup{}
	processGroup := &sync.WaitGroup{}

	// Add number of routines to the wait group, so the group know when it is done
	processGroup.Add(routines)
	logger.LogTrace(className, "Wait group created. Added routines number to it")

	// Create goroutines (with the go key word) and associate the process data function. It means, that when the channel
	// dataReceive starts to receive data, some of the goroutine took the data, process it, and sent the result on the
	// results channel
	for w := 1; w <= routines; w++ {
		logger.LogTrace(className, fmt.Sprintf("Creating worker %d", w))
		go e.processData(w, dataReceiveCh, resultsCh, processGroup)
	}

	// Declare results variables
	counter := new(int)          // int pointer
	successCounter := new(int64) // int64 pointer
	failCounter := new(int64)    // int64 pointer

	*counter = 0
	*successCounter = 0
	*failCounter = 0

	// Create a single goroutine for starting to write Data, whenever it starts to receive results in the channel
	generalGroup.Add(1)
	go e.writeData(resultsCh, generalGroup, flushNumber, counter, successCounter, failCounter)

	// Create a single goroutine, for start to read data from input file, and send it to the dataReceive channel
	generalGroup.Add(1)
	go e.readData(dataReceiveCh, generalGroup)

	// Wait until the execution finish, so it can close the last channel
	e.waitAndCloseChannel(resultsCh, generalGroup, processGroup)

	// After the group.Wait() its Done, here is finished the parallel section
	e.logResult(counter, successCounter, failCounter)
}

/*+
This function, receive a send only channel for send the data read, and then close the channel.
*/
func (e *ParallelExecutor) readData(dataReceiveCh chan<- []string, generalGroup *sync.WaitGroup) {
	logger.LogTrace(className, "Starting to consume data")
	for {
		line, err := e.dataLoader.ReadNextLine()
		if err == io.EOF {
			logger.LogDebug(className, "End of file reached for input file")
			break
		} else if err != nil {
			logger.LogError(className, fmt.Sprintf("Error reading line: %s", err.Error()))
		}

		// Send read data to channel
		dataReceiveCh <- line
	}
	logger.LogTrace(className, "Finished to consume data")
	generalGroup.Done()
	// Close the channel so the for in processDataFunction can finish and all the goroutines for process data can be done
	close(dataReceiveCh)
	logger.LogTrace(className, "Channel of reading data closed")
}

/*+
This function, receive data from a receive only channel (record), process the data with the processor, and then send
the result to the send only channel (result).
Lastly, when it knows that there is no more data for process (with the for range in the channel), it marks that the
goroutine is done in the wait group
*/
func (e *ParallelExecutor) processData(id int, dataReceiveCh <-chan []string, resultsCh chan<- []string, processGroup *sync.WaitGroup) {
	logger.LogTrace(className, fmt.Sprintf("Worker with ID %d started", id))

	for record := range dataReceiveCh {
		logger.LogTrace(className, fmt.Sprintf("Worker %d processing record: %s", id, record[0]))

		success, description := e.processor.Process(record)
		currentResult := append(record, strconv.FormatBool(success), description)

		logger.LogTrace(className, fmt.Sprintf("Worker %d finished record: %s", id, record[0]))

		resultsCh <- currentResult
	}

	logger.LogTrace(className, fmt.Sprintf("Worker with ID %d done", id))
	processGroup.Done()
}

/*+
This function, receive a receive only channel for read the results, and then write in the output file. Lastly,
make the goroutine done in the wait group
*/
func (e *ParallelExecutor) writeData(resultsCh <-chan []string, generalGroup *sync.WaitGroup, flushNumber int, counter *int, successCounter, failCounter *int64) {

	logger.LogTrace(className, "Starting to wait for results")
	for result := range resultsCh {
		*counter++

		logger.LogTrace(className, fmt.Sprintf("Starting to process result %d", *counter))

		if *counter%flushNumber == 0 {
			e.resultWriter.FlushWriter()
		}

		success, err := strconv.ParseBool(result[len(result)-2])
		if err != nil {
			logger.LogError(className, fmt.Sprintf("Error parsing success flag from result %s. Error: %s", result, err.Error()))
			continue
		}

		if success {
			*successCounter++
		} else {
			*failCounter++
		}

		_ = e.resultWriter.WriteResult(result, success)

		logger.LogInfo(logger.ColorYellow, className, fmt.Sprintf("Processed [success:%t] - %d", success, *counter))
	}

	logger.LogTrace(className, "Finished to write results")
	generalGroup.Done()
}

func (e *ParallelExecutor) waitAndCloseChannel(resultsCh chan []string, generalGroup *sync.WaitGroup, processGroup *sync.WaitGroup) {
	logger.LogTrace(className, "Starting to wait")
	processGroup.Wait()
	// Close the channel so the for in writeData can finish then it can sure that all the results are written
	close(resultsCh)

	// Wait for writerData to finish (and read data to, but it will be always done because the previous group finish when
	// the readData channel is closed)
	generalGroup.Wait()
	logger.LogTrace(className, "Results Closed")
}

func (e *ParallelExecutor) logResult(counter *int, successCounter, failCounter *int64) {
	logger.LogInfo(logger.ColorBlack, className, fmt.Sprintf("-------------------------------------------------------------------------------------"))
	logger.LogInfo(logger.ColorGreen, className, fmt.Sprintf("Successful records - %d", *successCounter))
	logger.LogInfo(logger.ColorRed, className, fmt.Sprintf("Failed records - %d", *failCounter))
	logger.LogInfo(logger.ColorTeal, className, fmt.Sprintf("Total records - %d", *counter))
}
