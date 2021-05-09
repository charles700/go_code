信号监听
----

### Notify 
```golang
func Notify(c chan<- os.Signal, sig ...os.Signal)
```
- 第一个参数表示接收信号的管道
- 第二个及后面的参数表示设置要监听的信号，如果不设置表示监听所有的信号。