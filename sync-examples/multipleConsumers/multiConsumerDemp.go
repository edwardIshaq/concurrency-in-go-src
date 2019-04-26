package main

import (
	"fmt"
	"sync"
)

func main() {
	singleProducer := producer()
	var wg sync.WaitGroup
	wg.Add(5)
	for cIdx := 0; cIdx < 5; cIdx++ {
		go func(cIdx int, producer <-chan string) {
			defer wg.Done()
			name := fmt.Sprintf("Worker_%d", cIdx)
			consumer(name, producer)
		}(cIdx, singleProducer)
	}
	wg.Wait()
}

func producer() <-chan string {
	valueStream := make(chan string)
	go func() {
		defer close(valueStream)
		for i := 0; i < 100; i++ {
			value := fmt.Sprintf("File_%d.txt", i)
			valueStream <- value
		}
	}()
	return valueStream
}

func consumer(name string, task <-chan string) {
	for {
		file, ok := <-task
		if !ok {
			fmt.Printf("\nconsumer: %s is done", name)
			return
		}
		fmt.Printf("\nconsumer: %s is handling file: %s", name, file)
	}
}
