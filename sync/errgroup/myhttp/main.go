package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, r.URL.Path)
	})

	server := &http.Server{Addr: ":8080"}

	g.Go(func() error {
		err := server.ListenAndServe()
		return err
	})

	// 监听 整个 errgroup 退出
	g.Go(func() error {
		<-ctx.Done()
		log.Println("err group 退出: 关闭 http 服务")
		return server.Shutdown(ctx)
	})

	// 监听 linux 信号 ctrl +c ， linux 信号退出时
	g.Go(func() error {
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

		s := <-signalCh
		log.Println("linux 退出信号'")
		return errors.Errorf("linux 信号退出: %v", s)
	})

	err := g.Wait()
	fmt.Printf("程序退出-- %v \n", err)
}
