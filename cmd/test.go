package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	var x int

	for {
		fmt.Scan(&x)
		switch x {
		case 1:
			cancel()
		case 0:
			ctx, cancel = context.WithCancel(context.TODO())
			go Listen(ctx)
		default:
			continue
		}

	}

}

func Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Println(".")
			time.Sleep(time.Second)
		}
	}
}
