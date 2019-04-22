package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	demoRepeat()
}

func timeoutDemo() {
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("Time out")
	}
}

func demoSelectTimeAfter() {
	fmt.Println("Start")
	startTime := time.Now()
loop:
	for {
		select {
		case <-time.After(time.Second):
			break loop
		}
	}
	fmt.Println("end")
	fmt.Printf("Took %v", time.Since(startTime))
}

func demoDoneGoroutine() {
	fmt.Println("Start")
	done := make(chan interface{})

	var wg sync.WaitGroup
	wg.Add(1)

	go func(done <-chan interface{}) {
		defer func() {
			fmt.Println("gorutine defer")
			wg.Done()
		}()

		for {
			select {
			case <-done:
				fmt.Println("Received done")
				return
			}
		}
	}(done)

	go func() {
		time.Sleep(time.Second)
		done <- struct{}{}
	}()

	wg.Wait()
	fmt.Println("End")
}

func demoRepeat() {
	fmt.Println("start")
	startTime := time.Now()
	done := make(chan interface{})
	repeatChan := repeat(done, "spin")

loop:
	for {
		select {
		case <-time.After(1 * time.Millisecond):
			fmt.Println("Timed out")
			close(done)
			break loop

		case value := <-repeatChan:
			fmt.Println(value)
		}
	}

	fmt.Println("end")
	fmt.Printf("took %v\n", time.Since(startTime))
}

func repeat(done <-chan interface{}, value interface{}) <-chan interface{} {
	output := make(chan interface{})
	go func() {
		defer close(output)
		for {
			select {
			case <-done:
				fmt.Println("Okay Thats enough")
				return
			case output <- value:
			}
		}
	}()
	return output
}

func stopWithSecondGoRoutine() {
	done := make(chan interface{})
	defer close(done)
	fmt.Println("start")
	repeatChan := repeat(done, "spin")

	go func() {
		<-time.After(time.Second)
		done <- struct{}{}
	}()

	for value := range repeatChan {
		fmt.Println(value)
	}
	fmt.Println("end")
}

func sleep(done <-chan struct{}, duration time.Duration, input <-chan struct{}) <-chan struct{} {
	output := make(chan struct{})
	go func() {
		defer close(output)
		for {
			input := input
			var val struct{}
			select {
			case <-done:
				return
			case val = <-input:
				output <- val
				input = nil
			case <-time.After(duration):
				select {
				case <-done:
					return
				case output <- val:
					continue
				}
			}
		}
	}()
	return output
}
