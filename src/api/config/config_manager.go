package config

import (
	"fmt"
	"os"

	"github.com/hmalatini/request_executor_script/src/api/config/models"
	"github.com/hmalatini/request_executor_script/src/api/logger"
	"gopkg.in/yaml.v2"
)

const className = "ConfigManager"

var cfg models.Config

func InitConfig(cfgPath string) error {
	if cfgPath == "" {
		cfgPath = "./src/api/config/config.yml"
	}

	f, err := os.Open(cfgPath)
	if err != nil {
		logger.LogError(className, fmt.Sprintf("Can not open configuration file: %s", err.Error()))
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		logger.LogError(className, fmt.Sprintf("Can not decode yml file: %s", err.Error()))
		return err
	}

	return nil
}

func GetConfig() models.Config {
	return cfg
}

func PrintConfig() {
	logger.LogInfo(logger.ColorMagenta, className, "-------------------------------------------------------------------------------------")
	logger.LogInfo(logger.ColorPurple, className, "STARTING SCRIPT")
	logger.LogInfo(logger.ColorMagenta, className, "Author: Hernan Malatini")
	logger.LogInfo(logger.ColorPurple, className, "EXECUTOR")
	logger.LogInfo(logger.ColorBlack, className, fmt.Sprintf("Routines number: %d", cfg.Executor.Routines))
	logger.LogInfo(logger.ColorPurple, className, "REQUEST")
	logger.LogInfo(logger.ColorBlack, className, fmt.Sprintf("Method: %s", cfg.Processor.Request.Method))
	logger.LogInfo(logger.ColorBlack, className, fmt.Sprintf("URL: %s%s", cfg.Processor.Request.BaseUrl, cfg.Processor.Request.Uri))
	logger.LogInfo(logger.ColorBlack, className, fmt.Sprintf("Body: %s", cfg.Processor.Request.Body))
	logger.LogInfo(logger.ColorBlack, className, fmt.Sprintf("Headers: %v", cfg.Processor.Request.Headers))
	logger.LogInfo(logger.ColorMagenta, className, "-------------------------------------------------------------------------------------")
}
