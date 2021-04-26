package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)
	go func() {
		time.Sleep(time.Second * 2)
		ch <- 1
	}()

	select {
	case a := <-ch:
		fmt.Println(a, "ch111")
	default:
		fmt.Println("default")
	}
}
