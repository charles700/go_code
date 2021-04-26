package main

import (
	"fmt"
	"runtime"
	"time"
)

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
