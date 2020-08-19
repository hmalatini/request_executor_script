package main

import (
	"flag"
	"fmt"
	"github.com/hmalatini/request_executor_script/src/api/config/models"

	"github.com/hmalatini/request_executor_script/src/api/config"
	dataLoaderPkg "github.com/hmalatini/request_executor_script/src/api/data_loader"
	"github.com/hmalatini/request_executor_script/src/api/executor"
	"github.com/hmalatini/request_executor_script/src/api/logger"
	"github.com/hmalatini/request_executor_script/src/api/processor"
	resultWriterPkg "github.com/hmalatini/request_executor_script/src/api/result_writer"
)

const className = "main"

func main() {
	configPath := flag.String("config", "", "The path of the config file")
	requestExample := flag.Bool("example", false, "Takes first line of input and print how the "+
		"request could be executed, without execute anything")

	flag.Parse()

	err := config.InitConfig(*configPath)
	if err != nil {
		return
	}

	cfg := config.GetConfig()
	if !*requestExample {
		config.PrintConfig()
	}

	logger.SetCurrentLogLevel(cfg.Logger.Level)

	dataLoader := dataLoaderPkg.NewCsvDataLoader()
	err = dataLoader.InitConnection()
	if err != nil {
		return
	}

	resultWriter := resultWriterPkg.NewCsvResultWriter()
	err = resultWriter.InitConnection()
	if err != nil {
		return
	}

	resultWriter.WriteHeader(dataLoader.GetHeaders())

	requestProcessor := processor.NewRequestProcessor()
	requestProcessor.InitProcessor(
		cfg.Processor.Request.BaseUrl,
		cfg.Processor.Request.Uri,
		cfg.Processor.Request.Body,
		cfg.Processor.Request.Headers,
		cfg.Processor.Request.Method)

	if *requestExample {
		printExampleRequest(cfg, *dataLoader, *requestProcessor)
		return
	}

	parallelExecutor := executor.NewParallelExecutor(*dataLoader, *requestProcessor, *resultWriter)
	parallelExecutor.Execute(cfg.Executor.Routines, cfg.Executor.Flush)

	resultWriter.CloseConnection()
}

func printExampleRequest(cfg models.Config, dataLoader dataLoaderPkg.CsvDataLoader, processor processor.RequestProcessor) {
	logger.LogInfo(logger.ColorGreen, className, "Example Request:")
	record, err := dataLoader.ReadNextLine()
	if err != nil {
		logger.LogError(className, "Error reading first line")
		return
	}

	uri, err := processor.GetParsedUri(record)
	if err != nil {
		return
	}

	body, err := processor.GetParsedBody(record)
	if err != nil {
		return
	}

	logger.LogInfo(logger.ColorBlack, className, fmt.Sprintf("Method: %s", cfg.Processor.Request.Method))
	logger.LogInfo(logger.ColorBlack, className, fmt.Sprintf("URL: %s%s", cfg.Processor.Request.BaseUrl, uri))
	logger.LogInfo(logger.ColorBlack, className, fmt.Sprintf("Body: %s", body))
	logger.LogInfo(logger.ColorBlack, className, fmt.Sprintf("Headers: %v", cfg.Processor.Request.Headers))
}
