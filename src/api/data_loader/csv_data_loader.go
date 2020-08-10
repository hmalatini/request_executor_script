package data_loader

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/hmalatini/request_executor_script/src/api/config"
	"github.com/hmalatini/request_executor_script/src/api/logger"
)

const className = "CsvDataLoader"

var reader *csv.Reader
var header []string

func init() {
	reader = nil
	header = nil
}

type CsvDataLoader struct{}

func NewCsvDataLoader() *CsvDataLoader {
	return &CsvDataLoader{}
}

func (d *CsvDataLoader) InitConnection() error {
	cfg := config.GetConfig()

	inputFile, err := os.Open(cfg.Data.File)
	if err != nil {
		logger.LogError(className, fmt.Sprintf("Error opening input csv file: %s", err.Error()))
		return err
	}
	logger.LogDebug(className, fmt.Sprintf("File with path %s has opened correctly", cfg.Data.File))

	reader = csv.NewReader(bufio.NewReader(inputFile))
	if !cfg.Data.Headers {
		return nil
	}

	header, err = reader.Read()
	if err != nil {
		logger.LogError(className, fmt.Sprintf("Error opening the input csv file: %s", err.Error()))
		return err
	}

	logger.LogDebug(className, fmt.Sprintf("Header -> %s", header))

	return nil
}

func (d *CsvDataLoader) ReadNextLine() ([]string, error) {
	line, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			logger.LogDebug(className, "EOF reached")
			return nil, err
		}

		logger.LogError(className, fmt.Sprintf("Error reading csv: %s", err.Error()))
		return nil, err
	}

	logger.LogDebug(className, fmt.Sprintf("Line read -> %s", line))

	return line, nil
}

func (d *CsvDataLoader) CloseConnection() {
}

func (d *CsvDataLoader) GetHeaders() []string {
	return header
}
