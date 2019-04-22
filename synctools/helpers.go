package synctools

import "fmt"

// Repeat will repeat a value untill done
func Repeat(done <-chan interface{}, value interface{}) <-chan interface{} {
	valueFunc := func() interface{} {
		return value
	}
	return RepeatFunc(done, valueFunc)
}

// RepeatFunc will continue caling fun untill done
func RepeatFunc(done <-chan interface{}, fun func() interface{}) <-chan interface{} {
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

// OrDone keeps pulling from stream untill done
func OrDone(done, stream <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-stream:
				if ok == false {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}
