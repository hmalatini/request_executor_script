package main

import (
	"flag"

	"github.com/hmalatini/request_executor_script/src/api/config"
	dataLoaderPkg "github.com/hmalatini/request_executor_script/src/api/data_loader"
	"github.com/hmalatini/request_executor_script/src/api/executor"
	"github.com/hmalatini/request_executor_script/src/api/logger"
	"github.com/hmalatini/request_executor_script/src/api/processor"
	resultWriterPkg "github.com/hmalatini/request_executor_script/src/api/result_writer"
)

func main() {
	configPath := flag.String("config", "", "The path of the config file")

	flag.Parse()

	err := config.InitConfig(*configPath)
	if err != nil {
		return
	}

	cfg := config.GetConfig()
	config.PrintConfig()

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

	parallelExecutor := executor.NewParallelExecutor(*dataLoader, *requestProcessor, *resultWriter)
	parallelExecutor.Execute(cfg.Executor.Routines, cfg.Executor.Flush)

	resultWriter.CloseConnection()
}
