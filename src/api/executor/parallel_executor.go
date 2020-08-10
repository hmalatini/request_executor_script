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
	var successCounter int64
	var failCounter int64

	//Declare the jobs to work, it means, the data for make the job
	records := make(chan []string, flushNumber)
	//Declare the results
	results := make(chan []string, flushNumber)

	// Initialize the job for start to work.. It started blocked until the reader start to read data
	group := sync.WaitGroup{}
	group.Add(routines)
	for w := 1; w <= routines; w++ {
		logger.LogTrace(className, fmt.Sprintf("Creating worker %d", w))
		go e.worker(w, records, results, &group)
	}

	go func() {
		logger.LogTrace(className, "Starting to wait")
		group.Wait()
		close(results)
		logger.LogTrace(className, "Results Closed")
	}()

	logger.LogTrace(className, "Sending to parallel read")
	go e.parallelRead(records)

	count := 0
	logger.LogTrace(className, "Starting to wait for results")

	for result := range results {
		logger.LogTrace(className, fmt.Sprintf("Entre acá: Counter: %d", count))
		count++

		if count%flushNumber == 0 {
			e.resultWriter.FlushWriter()
		}

		success, err := strconv.ParseBool(result[len(result)-2])
		if err != nil {
			logger.LogError(className, fmt.Sprintf("Error parsing success flag from result %s. Error: %s", result, err.Error()))
		}

		if success {
			successCounter++
		} else {
			failCounter++
		}

		_ = e.resultWriter.WriteResult(result, success)

		logger.LogInfo(logger.ColorYellow, className, fmt.Sprintf("Processed [success:%t] - %d", success, count))
	}

	logger.LogInfo(logger.ColorBlack, className, fmt.Sprintf("-------------------------------------------------------------------------------------"))
	logger.LogInfo(logger.ColorGreen, className, fmt.Sprintf("Successful records - %d", successCounter))
	logger.LogInfo(logger.ColorRed, className, fmt.Sprintf("Failed records - %d", failCounter))
	logger.LogInfo(logger.ColorTeal, className, fmt.Sprintf("Total records - %d", count))
}

func (e *ParallelExecutor) worker(id int, records <-chan []string, results chan<- []string, group *sync.WaitGroup) {
	logger.LogTrace(className, fmt.Sprintf("Worker with ID %d started", id))
	defer func() {
		logger.LogTrace(className, fmt.Sprintf("Worker with ID %d done", id))
		group.Done()
	}()

	for record := range records {
		logger.LogTrace(className, fmt.Sprintf("Worker %d processing record: %s", id, record[0]))

		success, description := e.processor.Process(record)
		currentResult := append(record, strconv.FormatBool(success), description)

		logger.LogTrace(className, fmt.Sprintf("Worker %d finished record: %s", id, record[0]))

		results <- currentResult
	}
}

func (e *ParallelExecutor) parallelRead(records chan []string) {
	logger.LogTrace(className, "Starting reading file")
	for {
		line, err := e.dataLoader.ReadNextLine()
		if err == io.EOF {
			break
		} else if err != nil {
			logger.LogError(className, fmt.Sprintf("Error reading line: %s", err.Error()))
		}

		records <- line
	}
	logger.LogTrace(className, "Finished reading file")
	close(records)
}
