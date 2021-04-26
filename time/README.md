time
---- 

## time.NewTicker
创建一个定时器，返回 一个 ticker, 有以属性和方法：
- `ticker.C chan Time`  获取系统时间的通道
- `ticker.Stop()` 结束定时器

```golang
func main() {
	ticker := time.NewTicker(2 * time.Second)
	fmt.Println("[0] 当前时间为：", time.Now().Format("2006-01-02 15:04:05"))

	var count = 0

	go func() {
		for {
			count++

			// 从定时器中读取数据
			t := <-ticker.C
			fmt.Println("[1] 当前时间为：:", t.Format("2006-01-02 15:04:05"))

			if count >= 5 {
				ticker.Stop() // 结束定时器
				runtime.Goexit()
			}
		}
	}()

	for {
		time.Sleep(time.Second)
	}
}
```


## time.NewTimer
创建一个延时器，返回 一个 timer, 有以属性和方法：
- `timer.C chan Time`  获取系统时间的通道
- `timer.Stop()` 结束定时器
- `timer.Reset(duration)` 重置定时器

```golang

var wg sync.WaitGroup

func main() {
	timer := time.NewTimer(time.Second * 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			t := <-timer.C
			fmt.Println(t.Format("2006-01-02 15:04:05"))
			// 通过 Reset 方法 达到定时执行的效果
			timer.Reset(time.Second * 2)
		}
	}()

	wg.Wait()
}

```
