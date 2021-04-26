package main

import (
	"fmt"
	"sync"
	"time"
)

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
