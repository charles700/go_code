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

// Logger 日志结构体
type Logger struct {
	Level LogLevel
}

//
func NewLog(levelstr string) Logger {
	level, err := parseLogLevel(levelstr)
	if err != nil {
		panic(err)
	}
	return Logger{
		Level: level,
	}
}

func (l *Logger) enable(level LogLevel) bool {
	return level >= l.Level
}

func log(lv LogLevel, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	now := time.Now()
	funcName, fileName, lineNo := getInfo(3)
	// "[2006-01-02 15:04:05] [DEBUG] [文件名:函数名:行号] 日志信息"
	fmt.Printf("[%s] [%s] [%s:%s:%d] %s \n", now.Format(("2006-01-02 15:04:05")), getLogLevelString(lv), fileName, funcName, lineNo, msg)
}

func (l *Logger) Debug(format string, a ...interface{}) {
	if l.enable(DEBUG) {
		log(DEBUG, format, a...)
	}
}

func (l *Logger) Info(format string, a ...interface{}) {
	if l.enable(INFO) {
		log(INFO, format, a...)
	}
}

func (l *Logger) Fatal(format string, a ...interface{}) {
	if l.enable(FATAL) {
		log(FATAL, format, a...)
	}
}

func (l *Logger) Error(format string, a ...interface{}) {
	if l.enable(ERROR) {
		log(ERROR, format, a...)
	}
}
func (l *Logger) Trace(format string, a ...interface{}) {
	if l.enable(TRACE) {
		log(TRACE, format, a...)
	}
}
func (l *Logger) Warning(format string, a ...interface{}) {
	if l.enable(WARNING) {
		log(WARNING, format, a...)
	}
}
