package example

import (
	"fmt"
	"github.com/hmalatini/request_executor_script/src/api/config"
	dataLoaderPkg "github.com/hmalatini/request_executor_script/src/api/data_loader"
	"github.com/hmalatini/request_executor_script/src/api/logger"
	"github.com/hmalatini/request_executor_script/src/api/processor"
	resultWriterPkg "github.com/hmalatini/request_executor_script/src/api/result_writer"
)

const className = "ExampleExecutor"

type ExampleExecutor struct {
	dataLoader   dataLoaderPkg.CsvDataLoader
	processor    processor.RequestProcessor
	resultWriter resultWriterPkg.CsvResultWriter
}

func NewExampleExecutor(dataLoader dataLoaderPkg.CsvDataLoader,
	processor processor.RequestProcessor,
	resultWriter resultWriterPkg.CsvResultWriter) *ExampleExecutor {

	return &ExampleExecutor{
		dataLoader:   dataLoader,
		processor:    processor,
		resultWriter: resultWriter,
	}
}

func (e *ExampleExecutor) Execute() {
	cfg := config.GetConfig()

	logger.LogInfo(logger.ColorGreen, className, "Example Request:")
	record, err := e.dataLoader.ReadNextLine()
	if err != nil {
		logger.LogError(className, "Error reading first line")
		return
	}

	uri, err := e.processor.GetParsedUri(record)
	if err != nil {
		return
	}

	body, err := e.processor.GetParsedBody(record)
	if err != nil {
		return
	}

	logger.LogInfo(logger.ColorBlack, className, fmt.Sprintf("Method: %s", cfg.Processor.Request.Method))
	logger.LogInfo(logger.ColorBlack, className, fmt.Sprintf("URL: %s%s", cfg.Processor.Request.BaseUrl, uri))
	logger.LogInfo(logger.ColorBlack, className, fmt.Sprintf("Body: %s", body))
	logger.LogInfo(logger.ColorBlack, className, fmt.Sprintf("Headers: %v", cfg.Processor.Request.Headers))
}
