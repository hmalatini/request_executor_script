package executor

import (
	dataLoaderPkg "github.com/hmalatini/request_executor_script/src/api/data_loader"
	"github.com/hmalatini/request_executor_script/src/api/executor/example"
	"github.com/hmalatini/request_executor_script/src/api/executor/parallel"
	"github.com/hmalatini/request_executor_script/src/api/processor"
	resultWriterPkg "github.com/hmalatini/request_executor_script/src/api/result_writer"
	"strings"
)

type InterfaceExecutor interface {
	Execute()
}

type executorFactory struct{}

func NewExecutorFactory() *executorFactory {
	return &executorFactory{}
}

func (ef *executorFactory) GetExecutor(executorName string,
	dataLoader dataLoaderPkg.CsvDataLoader,
	processor processor.RequestProcessor,
	resultWriter resultWriterPkg.CsvResultWriter) InterfaceExecutor {
	executorName = strings.ToLower(executorName)

	if executorName == "parallel" {
		return parallel.NewParallelExecutor(dataLoader, processor, resultWriter)
	} else if executorName == "example" {
		return example.NewExampleExecutor(dataLoader, processor, resultWriter)
	}

	return nil
}
