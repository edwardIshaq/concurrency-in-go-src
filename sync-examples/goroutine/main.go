package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	// goroutineSize()
	lockStep()
}

func example1() {
	runtime.GOMAXPROCS(3)
	task := func(name string, step, n int) {
		fmt.Printf("• %s step: %d count: %d\n", name, step, n)
		counter := 0
		for x := 0; x <= n; x++ {
			fmt.Printf("%s: %d\n", name, counter)
			counter += step
		}
	}

	go task("A", 2, 100)
	go task("B", 3, 100)
	time.Sleep(1 * time.Second)
}

// lockStepTimer will run two goroutines that alternate and exit after 1 second
func lockStepTimer() {
	runtime.GOMAXPROCS(2)

	task := func(output chan int, count int, name string) {
		for current := range output {
			fmt.Printf("%s = %d\n", name, current)
			current++
			if current == count {
				close(output)
				return
			}
			output <- current
		}
	}

	output := make(chan int)
	go task(output, 50, "A")
	go task(output, 50, "B")
	output <- 0
	time.Sleep(1 * time.Second)
}

// lockStepTimer will run two goroutines that alternate and exit when WaitGroup is done
func lockStep() {
	runtime.GOMAXPROCS(1)

	task := func(output chan int, count int, wg *sync.WaitGroup, name string) {
		defer wg.Done()
		for current := range output {
			fmt.Printf("%s = %d\n", name, current)
			current++
			if current == count {
				close(output)
				return
			}
			output <- current
		}
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)
	output := make(chan int)
	go task(output, 50, wg, "A")
	go task(output, 50, wg, "B")
	output <- 0
	wg.Wait()
}

func goroutineSize() {
	memConsumed := func() uint64 {
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}

	var c <-chan interface{}
	var wg sync.WaitGroup
	noop := func() { wg.Done(); <-c }

	const numGoroutines = 1e4
	fmt.Printf("Running %.1f\n", numGoroutines)
	wg.Add(numGoroutines)
	before := memConsumed()
	for i := numGoroutines; i > 0; i-- {
		go noop()
	}
	wg.Wait()
	after := memConsumed()
	fmt.Printf("%.3fkb", float64(after-before)/numGoroutines/1000)
}
