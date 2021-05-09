package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	g, _ := errgroup.WithContext(context.Background())

	for i := 0; i < 10; i++ {
		time.Sleep(200)
		index := i
		g.Go(func() error {
			time.Sleep(1000)
			fmt.Println(index)
			return fmt.Errorf("%d", index)
		})
	}

	err := g.Wait()
	fmt.Println("err", err)
}
