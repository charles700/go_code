package mylogger

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"strings"
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
