go 并发
---

## go 关键字
- go关键字可以快速创建 goroutine 
- goroutine 与 线程
  - 操作系统调度线程在可用处理器上运行
  - Go runtime 调度 goroutine 在绑定到单个操作系统线程的逻辑处理器中运行（P）。
    - 即使使用这个单一的逻辑处理器和操作系统线程，也可以调度数十万 goroutine 以惊人的效率和性能并发运行

## 使用 goroutine 的最佳实践
1. 如果 goroutine A 依赖另一个 goroutine B 的结果，才能工作。那么就应该 让 A 处理所有的工作。
  - 不使用B, 就消除了将结果从 B 返回到 A 所需的大量状态跟踪和 chan 操作。更加简单。

2. 让调用方决定 是否使用 goroutine
  - 如果函数内容 启动了 goroutine , 则必须向调用方提供显式停止 该 goroutine 的方法。

3. 在知道goroutine 什么时候终止，或 可以控制 goroutine 终止的情况下，再开启 goroutine 

4. 需要创建大量的goroutine 时，应当

## 细节
1. 空的 select 将阻塞 代码执行，最后退出程序
   -  defer 语句不会执行
2. log.Fatal()  会调用 os.Exit, 无条件 终止程序 -- 不建议使用 
   - defer 语句也不会执行
  
