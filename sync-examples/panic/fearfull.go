package main

import "fmt"

func main() {
	defer handlePanic()
	fearfull()
	fmt.Println(">>>>>>End")
}

func handlePanic() {
	fmt.Printf("calm down, its okay to: %v", recover())
}

func fearfull() {
	panic("nooooo")
	fmt.Println("will not reach here")
}
