package main

import (
	"fmt"
	"time"
)

func repeat(done <-chan interface{}, value interface{}) <-chan interface{} {
	valueFunc := func() interface{} {
		return value
	}
	return repeatFunc(done, valueFunc)
}

func repeatFunc(done <-chan interface{}, fun func() interface{}) <-chan interface{} {
	output := make(chan interface{})
	go func() {
		defer func() {
			fmt.Println("Okay Thats enough")
			close(output)
		}()
		for {
			select {
			case <-done:
				return
			case output <- fun():
			}
		}
	}()
	return output
}

func main() {
	fmt.Println("start")
	startTime := time.Now()
	done := make(chan interface{})
	repeatChan := repeat(done, "doing work")
	timerChan := time.After(10 * time.Millisecond)
loop:
	for {
		select {
		case <-timerChan:
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

func buggyMain() {
	fmt.Println("start")
	startTime := time.Now()
	done := make(chan interface{})
	repeatChan := repeat(done, "spin")

loop:
	for {
		select {
		case <-time.After(10 * time.Millisecond):
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
