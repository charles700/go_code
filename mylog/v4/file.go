package mylogger

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"time"
)

var (
	// MaxSize 日志通道的大小
	MaxSize = 50000
)

type logMsg struct {
	level     LogLevel
	msg       string
	funcName  string
	line      int
	fileName  string
	timestamp string
}

// 往 文件里面写日志
type FileLogger struct {
	Level       LogLevel
	fileObj     *os.File
	errFileObj  *os.File
	filePath    string // 日志文件保存的路径
	fileName    string // 日志文件名称
	maxFileSize int64
	logChan     chan *logMsg
}

// 构造函数
func NewFileLogger(levelstr, filePath, fileName string, maxFileSize int64) *FileLogger {
	logLevel, err := parseLogLevel(levelstr)
	if err != nil {
		panic(err)
	}

	fl := &FileLogger{
		Level:       logLevel,
		filePath:    filePath,
		fileName:    fileName,
		maxFileSize: maxFileSize,
		logChan:     make(chan *logMsg, MaxSize),
	}
	err = fl.initFile() // 按照文件路径和文件名 将文件打开
	if err != nil {
		panic(err)
	}

	return fl
}

func (f *FileLogger) initFile() error {
	fullFilePath := path.Join(f.filePath, f.fileName)
	fileObj, err := os.OpenFile(fullFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Printf("open log file failed, err: %v \n", err)
		return err
	}

	errFileObj, err := os.OpenFile(fullFilePath+".err", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open error log file failed, err: %v \n", err)
		return err
	}

	f.fileObj = fileObj
	f.errFileObj = errFileObj

	// 开启一个后台的goroutine 来记录日志
	// 多个 goroutine 写文件可能会报，会产生data race, 导致 日志丢失
	go f.writeLogBackground()
	return nil
}

// func (f *FileLogger) Close() {
// 	f.fileObj.Close()
// 	f.errFileObj.Close()
// }

func (l *FileLogger) enable(level LogLevel) bool {
	return level >= l.Level
}

func (l *FileLogger) checkSize(file *os.File) bool {
	fileInfo, err := file.Stat()

	if err != nil {
		fmt.Printf("get file info failed, err=%v \n", err)
		return false
	}
	return fileInfo.Size() >= l.maxFileSize
}

// checkTime  按照秒来切割 -- 没10秒切割一次
func (l *FileLogger) checkTime() bool {
	// 按秒来切割日志
	second := time.Now().Format("05")
	fmt.Printf("%s \n", time.Now().Format("2006-01-02 15:04:05"))
	i, err := strconv.Atoi(second)
	if err != nil {
		fmt.Printf("%v", err)
		return false
	}
	return i%10 == 0
}

func (l *FileLogger) splitFile(fileObj *os.File) (*os.File, error) {
	// 需要切割
	// 1. rename 一下 做备份  xx.log --> xx.log.bak20210304
	fileInfo, err := fileObj.Stat()
	nowStr := time.Now().Format("20060102150405")
	if err != nil {
		fmt.Printf("get file info failed, err=%v \n", err)
		return nil, err
	}
	logFilePath := path.Join(l.filePath, fileInfo.Name())
	newFileName := fmt.Sprintf("%s.bak%s", logFilePath, nowStr)
	os.Rename(logFilePath, newFileName)

	// 2. 关闭当前的日志文件
	fileObj.Close()

	// 3. 打开一个新的文件
	newFileObj, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("open new log file failed: %v", err)
		return nil, err
	}
	// 4. 将打开的新文件 赋值 给 l.fileobj
	return newFileObj, nil
}

func (l *FileLogger) writeLogBackground() {

	// 死循环 从 通道中取值
	for {

		// 检查大小是否需要切割
		if l.checkSize(l.fileObj) {
			// if l.checkTime() {
			newFileObj, _ := l.splitFile(l.fileObj)
			l.fileObj = newFileObj
		}

		// 1. 从通道中取值
		select {
		case logTmp := <-l.logChan:

			logInfo := fmt.Sprintf("[%s] [%s] [%s:%s:%d] %s \n", logTmp.timestamp, getLogLevelString(logTmp.level), logTmp.fileName, logTmp.funcName, logTmp.line, logTmp.msg)

			// "[2006-01-02 15:04:05] [DEBUG] [文件名:函数名:行号] 日志信息"
			fmt.Fprintf(l.fileObj, logInfo)

			if logTmp.level >= ERROR {
				if l.checkSize(l.errFileObj) {
					// if l.checkTime() {
					newFileObj, _ := l.splitFile(l.errFileObj)
					l.errFileObj = newFileObj
				}
				// 高于错误级别的日志多记录一次
				fmt.Fprintf(l.errFileObj, logInfo)
			}
		default:
			// 取到日志先 休息 500 毫秒，让出 CPU
			time.Sleep(time.Millisecond * 500)
		}

	}

}

func (l *FileLogger) log(lv LogLevel, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	now := time.Now()
	funcName, fileName, lineNo := getInfo(3)

	// 先把日志发送到 通道中
	// 1. 构造一个 logMsg 对象
	logTmp := &logMsg{
		level:     lv,
		msg:       msg,
		funcName:  funcName,
		line:      lineNo,
		fileName:  fileName,
		timestamp: now.Format("2006-01-02 15:04:05"),
	}

	// 2. 写入到 logChan 中去
	select {
	case l.logChan <- logTmp:
	default:
		// 当chan 满时，扔掉日志， 防止出现极端情况下出现阻塞
	}
}

func (l *FileLogger) Debug(format string, a ...interface{}) {
	if l.enable(DEBUG) {
		l.log(DEBUG, format, a...)
	}
}

func (l *FileLogger) Info(format string, a ...interface{}) {
	if l.enable(INFO) {
		l.log(INFO, format, a...)
	}
}

func (l *FileLogger) Fatal(format string, a ...interface{}) {
	if l.enable(FATAL) {
		l.log(FATAL, format, a...)
	}
}

func (l *FileLogger) Error(format string, a ...interface{}) {
	if l.enable(ERROR) {
		l.log(ERROR, format, a...)
	}
}
func (l *FileLogger) Trace(format string, a ...interface{}) {
	if l.enable(TRACE) {
		l.log(TRACE, format, a...)
	}
}
func (l *FileLogger) Warning(format string, a ...interface{}) {
	if l.enable(WARNING) {
		l.log(WARNING, format, a...)
	}
}
