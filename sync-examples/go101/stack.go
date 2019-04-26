package main

import (
	"fmt"
	"strconv"
)

type stack interface {
	Push(in []int, item int) (out []int)
	Top(in []int) (top int)
	Pop(in []int) (out []int)
	PopAndTop(in []int) (top int, out []int)
}

// Push to stack
func Push(in []int, item int) (out []int) {
	return append(in, item)
}

// Top of the stack
func Top(in []int) (top int) {
	if len(in) == 0 {
		panic("Cant top an empty stack")
	}
	return in[len(in)-1]
}

// Pop from stack
func Pop(in []int) (out []int) {
	lastIdx := len(in) - 1
	if lastIdx < 0 {
		lastIdx = 0
	}
	return in[0:lastIdx]
}

// PopAndTop Pop and return the element with the new stack
func PopAndTop(in []int) (top int, out []int) {
	top = Top(in)
	out = Pop(in)
	return
}

func testStack() {
	var myStack []int
	myStack = Push(myStack, 5)
	myStack = Push(myStack, 15)
	myStack = Push(myStack, 20)

	fmt.Println(myStack)
	fmt.Println(Top(myStack))

	myStack = Pop(myStack)
	fmt.Println(myStack)
	myStack = Pop(myStack)
	fmt.Println(myStack)
	top, myStack := PopAndTop(myStack)
	fmt.Println(top, myStack)
}

func main() {
	runCalc([]string{"1", "2", "+", "4", "*"}, 12)
	runCalc([]string{"4", "2", "/"}, 2)
	runCalc([]string{"10", "2", "5", "*", "+"}, 20)
}

func runCalc(expr []string, expectedValue int) {
	ans := postfixCalc(expr)
	fmt.Println(ans, ans == expectedValue)
}

// http://interactivepython.org/lpomz/courselib/static/pythonds/BasicDS/InfixPrefixandPostfixExpressions.html
func postfixCalc(expr []string) int {
	var stack = make([]int, 0, len(expr))
	var lhs, rhs int
	for _, val := range expr {
		num, err := strconv.Atoi(val)
		if err == nil { //found a number
			stack = Push(stack, num)
			continue
		}

		if len(stack) < 2 {
			panic("Expecting two operands per operator")
		}
		rhs, stack = PopAndTop(stack)
		lhs, stack = PopAndTop(stack)
		stack = Push(stack, strToFunc(val)(lhs, rhs))
	}
	return Top(stack)
}

type mathOp func(int, int) int

var (
	opMap = map[string]mathOp{
		"+": func(a, b int) int { return a + b },
		"-": func(a, b int) int { return a - b },
		"*": func(a, b int) int { return a * b },
		"/": func(a, b int) int { return a / b },
	}
)

func strToFunc(op string) mathOp {
	operator, ok := opMap[op]
	if !ok {
		panic(fmt.Sprintf("Unknown operator: %s", op))
	}
	return operator
}

func parseAsNumber(v string) *int {
	num, err := strconv.Atoi(v)
	if err != nil {
		return nil
	}
	return &num
}
