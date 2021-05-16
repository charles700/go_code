time
---- 

## 时间类型
```golang
func timeDemo() {
	now := time.Now() //获取当前时间
	fmt.Printf("current time:%v\n", now)

	year := now.Year()     //年
	month := now.Month()   //月
	day := now.Day()       //日
	hour := now.Hour()     //小时
	minute := now.Minute() //分钟
	second := now.Second() //秒
	fmt.Printf("%d-%02d-%02d %02d:%02d:%02d\n", year, month, day, hour, minute, second)
}
```


## 时间戳
时间戳是自1970年1月1日（08:00:00GMT）至当前时间的总毫秒数。它也被称为Unix时间戳（UnixTimestamp）。
```golang
func timestampDemo() {
	now := time.Now()            //获取当前时间
	timestamp1 := now.Unix()     //时间戳
	timestamp2 := now.UnixNano() //纳秒时间戳
	fmt.Printf("current timestamp1:%v\n", timestamp1)
	fmt.Printf("current timestamp2:%v\n", timestamp2)
}
```
### 将时间戳转为时间格式
使用time.Unix()函数可以将时间戳转为时间格式。
```golang
func timestampDemo2(timestamp int64) {
	timeObj := time.Unix(timestamp, 0) //将时间戳转为时间格式
	fmt.Println(timeObj)
	year := timeObj.Year()     //年
	month := timeObj.Month()   //月
	day := timeObj.Day()       //日
	hour := timeObj.Hour()     //小时
	minute := timeObj.Minute() //分钟
	second := timeObj.Second() //秒
	fmt.Printf("%d-%02d-%02d %02d:%02d:%02d\n", year, month, day, hour, minute, second)
}
```

## 时间间隔 
time.Duration是time包定义的一个类型，它代表两个时间点之间经过的时间，以纳秒为单位。time.Duration表示一段时间间隔，可表示的最长时间段大约290年。
```golang
const (
    Nanosecond  Duration = 1
    Microsecond          = 1000 * Nanosecond
    Millisecond          = 1000 * Microsecond
    Second               = 1000 * Millisecond
    Minute               = 60 * Second
    Hour                 = 60 * Minute
)
```

## 时间操作
```golang
// Add 操作
func (t Time) Add(d Duration) Time

now := time.Now()
later := now.Add(time.Hour) // 当前时间加1小时后的时间
fmt.Println(later)

// Sub 差值 - t-u
func (t Time) Sub(u Time) Duration

// 是否相等 - t == u
func (t Time) Equal(u Time) bool

// Before - 如果t代表的时间点在u之前，返回真；否则返回假。
func (t Time) Before(u Time) bool

// After - 如果t代表的时间点在u之后，返回真；否则返回假。
func (t Time) After(u Time) bool
```



## 定时器 -- time.NewTicker
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


## 定时器 --  time.NewTimer
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
