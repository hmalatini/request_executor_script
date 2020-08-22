package result_writer

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/hmalatini/request_executor_script/src/api/config"
	"github.com/hmalatini/request_executor_script/src/api/logger"
)

const className = "CsvResultWriter"

var writer *csv.Writer

type CsvResultWriter struct{}

func NewCsvResultWriter() *CsvResultWriter {
	return &CsvResultWriter{}
}

func (r *CsvResultWriter) InitConnection() error {
	cfg := config.GetConfig()

	var resultFile *os.File
	var err error

	resultFile, err = getFileOrCreate(cfg.Result.File, cfg.Result.Append)
	if err != nil {
		return err
	}

	writer = csv.NewWriter(resultFile)
	logger.LogDebug(className, "Writer Init Successfully")

	return nil
}

func (r *CsvResultWriter) WriteHeader(headers []string) error {
	if headers == nil {
		logger.LogDebug(className, "No writing headers because are empty")
		return nil
	}

	headers = append(headers, "Success", "Description")
	err := writer.Write(headers)
	if err != nil {
		logger.LogError(className, "Error writing headers")
		return err
	}

	logger.LogDebug(className, "Headers have written successfully")

	return nil
}

func (r *CsvResultWriter) WriteResult(result []string, success bool) error {

	if !r.hasToWrite(success) {
		logger.LogDebug(className, fmt.Sprintf("Discarted %s record", strconv.FormatBool(success)))
		return nil
	}

	err := writer.Write(result)
	if err != nil {
		logger.LogError(className, fmt.Sprintf("Error writing line %s", result))
		return err
	}

	logger.LogDebug(className, fmt.Sprintf("Record %s written successfully", result))

	return nil
}

func (r *CsvResultWriter) FlushWriter() {
	r.CloseConnection()
}

func (r *CsvResultWriter) CloseConnection() {
	writer.Flush()
	logger.LogDebug(className, "Writer Flushed")
}

func (r *CsvResultWriter) hasToWrite(success bool) bool {
	cfg := config.GetConfig()

	return (success && cfg.Result.WriteSuccess) || (!success && cfg.Result.WriteFails)
}
