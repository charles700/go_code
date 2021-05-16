package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const addr = ":9527"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "just another http server...")
	})
	//使用默认路由创建 http server
	srv := http.Server{
		Addr:    addr,
		Handler: http.DefaultServeMux,
	}

	//使用WaitGroup同步Goroutine
	var wg sync.WaitGroup
	exit := make(chan os.Signal)

	//监听 Ctrl+C 信号
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-exit

		wg.Add(1)
		//使用context控制srv.Shutdown的超时时间
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		fmt.Println("111")
		err := srv.Shutdown(ctx)
		fmt.Println("222")

		if err != nil {
			fmt.Println(err)
		}

		wg.Done()
	}()

	fmt.Println("listening at " + addr)
	err := srv.ListenAndServe()

	fmt.Println("waiting for the remaining connections to finish...")

	wg.Wait()

	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
	fmt.Println("gracefully shutdown the http server...")
}
