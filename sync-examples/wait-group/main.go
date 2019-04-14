package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	waitOrder()
}

func waitOrder() {
	runtime.GOMAXPROCS(2)
	var wg sync.WaitGroup
	for _, salutation := range []string{"Hello", "Yo yo yo", "Whats up?"} {
		wg.Add(1)
		go func(salutation string) {
			defer wg.Done()
			fmt.Println(salutation)
		}(salutation)
	}
	wg.Wait()
}
