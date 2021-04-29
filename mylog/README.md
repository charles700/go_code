日志库
---
## 需求分析
1. 支持往不同的地方输出日志
  - 终端
  - 文件
  - kafka
2. 日志分级别
  - Debug
  - Trace
  - Info
  - Warning
  - Error
  - Fatal
3. 支持开关
4. 完整日志信息: 时间、行号、文件名、日志级别、日志信息
5. 日志文件要切割

### v1 版本功能
> v1/console.go

- 终端输出
  - 调用 `fmt.Printf(format, msg)` 输出日志
- 日志分级别
  - 定义 `const`相关常量, 利用 `iota` 来分级别
- 支持开关
  - 增加 `func (c *ConsoleLogger) enable(level LogLevel) bool` 来控制级别开关
- 完整日志信息
  - 使用 `pc, file, lineNo, ok :=runtime.Caller(skip)` 获取当前执行的详细信息
  - 使用 `runtime.FuncForPC(pc).Name()` 获取执行的函数名
- 日志支持 格式化输出

### v2 版本功能 
> v2/file.go

- 文件输出
  - 使用  `fmt.Fprintf(file *os.File, format_msg string)`
- 日志文件切割
  - 按文件大小
    - 每次记录日志前，都判断下当前日志文件的大小， 
        - 新增 `checkSize()`函数获取日志文件大小, 利用 `file.Stat()` 返回的 `fileinfo.Size()`
    - 超过限制大小时，重命名文件。然后继续写日志  
        - 在初始化时 传入`maxFileSize` 来限制单个日志文件大小
        - 使用 `os.Rename` 来重命名旧的日志文件，然后继续向初始日志文件路径中写日志

  - 按日期分割


### v3 异步写日志 - goroutine



### v4 写入kafka



### v5 日志集成 ELK