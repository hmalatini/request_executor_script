package models

import "fmt"

type LogLevel int8

var logLevelStringToValue map[string]LogLevel
var logLevelValueToString map[LogLevel]string

const (
	LogLevelTrace = iota + 1
	LogLevelDebug
	LogLevelInfo
	LogLevelWarning
	LogLevelError
	LogLevelOff

	//Strings in config file
	LogLevelTraceString   = "TRACE"
	LogLevelDebugString   = "DEBUG"
	LogLevelInfoString    = "INFO"
	LogLevelWarningString = "WARNING"
	LogLevelErrorString   = "ERROR"
	LogLevelOffString     = "OFF"
)

func init() {
	logLevelStringToValue = map[string]LogLevel{
		LogLevelTraceString:   LogLevelTrace,
		LogLevelDebugString:   LogLevelDebug,
		LogLevelInfoString:    LogLevelInfo,
		LogLevelWarningString: LogLevelWarning,
		LogLevelErrorString:   LogLevelError,
		LogLevelOffString:     LogLevelOff,
	}

	logLevelValueToString = map[LogLevel]string{
		LogLevelTrace:   LogLevelTraceString,
		LogLevelDebug:   LogLevelDebugString,
		LogLevelInfo:    LogLevelInfoString,
		LogLevelWarning: LogLevelWarningString,
		LogLevelError:   LogLevelErrorString,
		LogLevelOff:     LogLevelOffString,
	}
}

func GetLogLevelValue(s string) (LogLevel, error) {
	result, ok := logLevelStringToValue[s]

	if !ok {
		return LogLevelInfo, fmt.Errorf("no value for %s", s)
	}

	return result, nil
}

func GetLogLevelString(level LogLevel) (string, error) {
	result, ok := logLevelValueToString[level]

	if !ok {
		return LogLevelInfoString, fmt.Errorf("no value for %d", level)
	}

	return result, nil
}
