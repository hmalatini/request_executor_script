package logger

import (
	"fmt"
	"github.com/hmalatini/request_executor_script/src/api/logger/models"
	"github.com/hmalatini/request_executor_script/src/api/utils"
)

var currentLevel models.LogLevel

const className string = "Logger"

func init() {
	currentLevel = models.LogLevelInfo
}

func SetCurrentLogLevel(level string) {
	var err error
	currentLevel, err = models.GetLogLevelValue(level)

	if err != nil {
		LogWarning(className, "Assign default log level info: "+err.Error())
	}

	assigned, _ := models.GetLogLevelString(currentLevel)
	LogDebug(className, fmt.Sprintf("%s assigned for Logger Level", assigned))
}

func LogTrace(class string, msg string) {
	if !shouldLog(models.LogLevelTrace) {
		return
	}

	fmt.Println(White(getMsgWithHeader(class, msg)))
}

func LogDebug(class string, msg string) {
	if !shouldLog(models.LogLevelDebug) {
		return
	}

	fmt.Println(Black(getMsgWithHeader(class, msg)))
}

func LogInfo(color LogColor, class string, msg string) {
	if !shouldLog(models.LogLevelInfo) {
		return
	}

	logMsg := getMsgWithHeader(class, msg)

	switch color {
	case ColorBlack:
		fmt.Println(Black(logMsg))
	case ColorGreen:
		fmt.Println(Green(logMsg))
	case ColorMagenta:
		fmt.Println(Magenta(logMsg))
	case ColorPurple:
		fmt.Println(Purple(logMsg))
	case ColorRed:
		fmt.Println(Red(logMsg))
	case ColorTeal:
		fmt.Println(Teal(logMsg))
	case ColorYellow:
		fmt.Println(Yellow(logMsg))
	default:
		fmt.Println(White(logMsg))
	}
}

func LogWarning(class string, msg string) {
	if !shouldLog(models.LogLevelWarning) {
		return
	}

	fmt.Println(Yellow(getMsgWithHeader(class, msg)))
}

func LogError(class string, msg string) {
	if !shouldLog(models.LogLevelError) {
		return
	}

	fmt.Println(Red(getMsgWithHeader(class, msg)))
}

func shouldLog(level models.LogLevel) bool {
	return currentLevel <= level
}

func getMsgWithHeader(class string, msg string) string {
	if currentLevel > models.LogLevelDebug {
		return msg
	}
	return fmt.Sprintf("%s - %s: %s", utils.GetStringTime(), class, msg)
}
