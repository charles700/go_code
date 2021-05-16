package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func worker(ctx context.Context) {
	time1 := time.Now()

LOOP1:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("done", time.Now().Sub(time1))
			break LOOP1
		default:
			continue
		}
	}

	wg.Done()
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	wg.Add(1)
	go worker(ctx)

	time.Sleep(time.Second * 3)
	// 3s 后调用 cancel 将会体检结束 ctx.Done
	// 不调用 cancel 的话，5秒后 ctx.Done 自动结束
	cancel()

	wg.Wait()

}
