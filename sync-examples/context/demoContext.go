package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream := generate(ctx)
	go consume(cancel, stream)
	<-ctx.Done()
}

func generate(ctx context.Context) <-chan int {
	output := make(chan int)
	go func(ctx context.Context) {
		defer close(output)
		counter := 0
		for {
			counter++
			if counter == 100 {
				return
			}
			select {
			case <-ctx.Done():
				return
			case output <- counter:
			}
		}
	}(ctx)
	return output
}

func consume(cancel context.CancelFunc, stream <-chan int) {
	defer cancel()
	for value := range stream {
		fmt.Println(value)
		time.Sleep(50 * time.Millisecond)
		if value == 200 {
			return
		}
	}
}
