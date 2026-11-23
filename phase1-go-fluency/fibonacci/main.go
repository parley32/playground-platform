package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run . <number>")
		fmt.Println("example: go run . 5")
		os.Exit(1)
	}

	input, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		fmt.Println("input is incorrect")
		os.Exit(1)
	}
	if input < 0 {
		fmt.Println("input must be greater than or equal to 0")
		os.Exit(1)
	}

	fmt.Printf("the fibonacci number at the %dth position is %d\n", input, getFibonacci(input))
}

func getFibonacci(input int64) int64 {
	x, y := int64(0), int64(1)
	for i := int64(0); i < input; i++ {
		x, y = y, x+y
	}
	return x
}
