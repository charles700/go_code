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


### v3 按日期分割
- 定义接口
- 按日期分割
```golang
// 没次遇到 10、20、30 等 秒时切割
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
```

### v4 异步写日志 - goroutine
要点：
- 在 `FileLogger` 结构体中加入一个通道
- 通道类型为一个[结构体的指针] -- 存储的数据量很小
  - 通道类型不能是string, 因为 string 是 int64 ，占 8个字节
- 在构造函数中就开启 `goroutine`从 通道中 读取日志,来写入文件
  - 使用 `select` ，并 默认增加 休眠时间，释放CPU
- 在统一的 `log` 方法中向通道写入日志数据


### v5 LogAgent

要点：
1. `kafka` 和 `zookeeper` 安装与启动
2.  使用 `hpcloud/tail` 模块, 监听文件新增内容，获取最新内容（读取最新行）
> go get github.com/hpcloud/tail
3. tail 接收到新日志后，使用 `sarama` 模块，连接kafka, 并发送给 kafka
4. 使用 `gopkg.in/ini.v1` 来加载配置文件


### v6 使用 etcd 作为日志手机的配置中心
目的：
1. `tiallog`根据 从 `etcd` 中获取的日志配置，监听日志文件的变化，把最新日志 发送给 `kafka`
2. 当一个新的配置文件在 `etcd` 中更新时，`etcd` 要 通知 `taillog` 模块做相应处理

#### 启动步骤
1. 服务配置文件 `conf.ini` 做映射
2. 启动`kafka` -- 根据配置
    - 配置`address`和`chan_size`
    - 初始化步骤
      - 与Kafka 服务建立链接
      - 创建一个生产者 `client` 和 消息通道 `logDataChan`
      - 监听通道，接收到日志消息，就发送到 Kafka 服务
3. 启动`etcd`  -- 根据配置
    - 配置`address`、`timeout`、`collect_log_key` 
    - 初始化步骤
      - 创建一个客户端，与 Etcd 服务建立连接
4. 从`etcd`中读取配置、监听 Etcd 服务中的日志配置数据变更
> etcd 中的初始日志配置数据有 `etcd_put/main.go` 事先塞入 
> 日志配置数据结构：`[]{topic, path}`
    - 使用 etcd 客户端 从 Etcd 服务中读取日志配置数据，并返回
    - 

5. 获取日志配置数据后，初始化 `taillog`服务
   - 初始化步骤
     - 每一条配置都启动一个 `TailTask` 任务
      - 使用 `tail` 模块 从`path`中读取信息，发送给 Kafka 的 `topic`
     - 创建一个通道`newConfChan`，监听日志配置数据的变更
        - 由其他人员更改了etcd 中的日志配置数据
        - etcd 监听到 变更，向通道喜欢滴信息
        - `taillog` 模块的 `newConfChan` 通道 接收到 变更后，再做相应处理


### v7 logTransfer
1. 从kafka 中把日志取出来，写入ES，使用Kibana 做可视化展示


### v8 系统监控
1. 系统监控
   - gopsutil 库 做系统监控、信息采集，写入 influxDB 中 使用 grafana 做展示
2. prometheus 监控、采集指标、保存， 使用 grafana 做展示


