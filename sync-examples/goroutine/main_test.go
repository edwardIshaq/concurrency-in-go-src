package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestBufferAppend(t *testing.T) {
	buff := []byte{}

	msg1 := []byte("Sender-")
	buff = append(buff, msg1...)

	msg2 := []byte("Receiver")
	buff = append(buff, msg2...)
	fmt.Println(string(buff))
	if string(buff) != "Sender-Receiver" {
		t.Errorf("failed")
	}
}

// BenchmarkContextSwitch test goroutines context switching
func BenchmarkContextSwitch(b *testing.B) {
	var wg sync.WaitGroup
	begin := make(chan struct{}, 1)
	c := make(chan struct{})
	var token struct{}

	sender := func() {
		defer wg.Done()
		<-begin

		for i := 0; i < b.N; i++ {
			c <- token
		}
	}

	receiver := func() {
		defer wg.Done()

		<-begin
		for i := 0; i < b.N; i++ {
			<-c
		}
	}

	wg.Add(2)
	go sender()
	go receiver()
	b.StartTimer()
	close(begin)
	wg.Wait()
}
