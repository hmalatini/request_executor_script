package result_writer

import (
	"fmt"
	"github.com/hmalatini/request_executor_script/src/api/logger"
	"github.com/hmalatini/request_executor_script/src/api/utils"
	"os"
)

const helperClassName = "ResultFileHelper"

func getFileOrCreate(path string, append bool) (*os.File, error) {
	var resultFile *os.File
	var err error

	if path != "" {
		if append {
			resultFile, err = os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
			logger.LogDebug(helperClassName, "Opened result file with append mode")
		} else {
			if fileExists(path) {
				err = os.Remove(path)
				if err != nil {
					logger.LogError(helperClassName, fmt.Sprintf("A result file already exists and cant be overwritten: %s", err.Error()))
					return nil, err
				}
				logger.LogDebug(helperClassName, fmt.Sprintf("File %s deleted", path))
			}
			resultFile, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
			logger.LogDebug(helperClassName, "Opened result file without append mode")
		}
		if err != nil {
			logger.LogError(helperClassName, fmt.Sprintf("Error opening result csv file: %s", err.Error()))
			return nil, err
		}
		logger.LogDebug(helperClassName, fmt.Sprintf("File with path %s has opened correctly", path))
	} else {
		resultFile, err = os.Create(fmt.Sprintf("./result-%s.csv", utils.GetStringTime()))
		if resultFile == nil {
			logger.LogError(helperClassName, "Error creating result csv file")
			return nil, err
		}
		logger.LogDebug(helperClassName, fmt.Sprintf("File with name %s has created correctly", resultFile.Name()))
	}

	return resultFile, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
