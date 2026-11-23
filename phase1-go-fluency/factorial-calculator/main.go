package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run . <number>")
		fmt.Println("example: go run . 6")
		os.Exit(1)
	}

	inputValue, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		fmt.Println("invalid input")
		os.Exit(1)
	}
	if inputValue < 0 {
		fmt.Println("input must be reater than or equal to 0")
		os.Exit(1)
	}
	if inputValue >= 21 {
		fmt.Println("Hold your horses there cow(girl/boy/person), this is an intro exercice, the number has to be strictly smaller than 21 or else the result will overflow")
		os.Exit(1)
	}
	factorial := int64(1)
	for value := int64(2); value <= inputValue; value++ {
		factorial *= value
	}
	fmt.Printf("The factorial of %d is %d\n", inputValue, factorial)
}
