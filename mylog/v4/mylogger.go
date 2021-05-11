package mylogger

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"strings"
)

// Logger 接口
type Logger interface {
	Debug(format string, a ...interface{})
	Info(format string, a ...interface{})
	Error(format string, a ...interface{})
	Warning(format string, a ...interface{})
	Fatal(format string, a ...interface{})
	Trace(format string, a ...interface{})
}

type LogLevel uint8

const (
	// 日志级别 自增
	UNKONWN LogLevel = iota
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)

func parseLogLevel(s string) (LogLevel, error) {
	s = strings.ToUpper(s)
	switch s {
	case "DEBUG":
		return DEBUG, nil
	case "TRACE":
		return TRACE, nil
	case "INFO":
		return INFO, nil
	case "WARNING":
		return WARNING, nil
	case "FATAL":
		return FATAL, nil
	case "ERROR":
		return ERROR, nil
	default:
		err := errors.New("无效的日志级别")
		return UNKONWN, err
	}
}
func getLogLevelString(lv LogLevel) string {
	switch lv {
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case FATAL:
		return "FATAL"
	case ERROR:
		return "ERROR"
	}
	return "DEBUG"
}

func getInfo(skip int) (funcName, fileName string, lineNo int) {
	pc, file, lineNo, ok := runtime.Caller(skip)

	if !ok {
		fmt.Println("runtime.Caller error")
		return
	}

	funcName = runtime.FuncForPC(pc).Name()
	fileName = path.Base(file)
	funcName = strings.Split(funcName, ".")[1]

	return
}
