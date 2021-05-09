package main

import (
	"fmt"
	"os"
	"os/signal"
	//"syscall"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c)
	s := <-c

	// 监听 ctrl +c 的退出信号
	if s == os.Interrupt {
		fmt.Println("get ctrl + c signal")
	} else {
		fmt.Println("get other signal")
	}
}
