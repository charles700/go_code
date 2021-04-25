package mylogger

// 往终端输出日志内容

import (
	"fmt"
	"time"
)

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

// ConsoleLogger 日志结构体
type ConsoleLogger struct {
	Level LogLevel
}

//
func NewLog(levelstr string) ConsoleLogger {
	level, err := parseLogLevel(levelstr)
	if err != nil {
		panic(err)
	}
	return ConsoleLogger{
		Level: level,
	}
}

func (c *ConsoleLogger) enable(level LogLevel) bool {
	return level >= c.Level
}

func log(lv LogLevel, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	now := time.Now()
	funcName, fileName, lineNo := getInfo(3)
	// "[2006-01-02 15:04:05] [DEBUG] [文件名:函数名:行号] 日志信息"
	fmt.Printf("[%s] [%s] [%s:%s:%d] %s \n", now.Format(("2006-01-02 15:04:05")), getLogLevelString(lv), fileName, funcName, lineNo, msg)
}

func (c *ConsoleLogger) Debug(format string, a ...interface{}) {
	if c.enable(DEBUG) {
		log(DEBUG, format, a...)
	}
}

func (c *ConsoleLogger) Info(format string, a ...interface{}) {
	if c.enable(INFO) {
		log(INFO, format, a...)
	}
}

func (c *ConsoleLogger) Fatal(format string, a ...interface{}) {
	if c.enable(FATAL) {
		log(FATAL, format, a...)
	}
}

func (c *ConsoleLogger) Error(format string, a ...interface{}) {
	if c.enable(ERROR) {
		log(ERROR, format, a...)
	}
}
func (c *ConsoleLogger) Trace(format string, a ...interface{}) {
	if c.enable(TRACE) {
		log(TRACE, format, a...)
	}
}
func (c *ConsoleLogger) Warning(format string, a ...interface{}) {
	if c.enable(WARNING) {
		log(WARNING, format, a...)
	}
}
