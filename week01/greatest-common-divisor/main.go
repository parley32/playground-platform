package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: go run . <number A> <number B>")
		fmt.Println("example: go run . 6 12")
		os.Exit(1)
	}

	inputA, errA := strconv.ParseInt(os.Args[1], 10, 64)
	inputB, errB := strconv.ParseInt(os.Args[2], 10, 64)
	if errA != nil || errB != nil {
		fmt.Println("invalid input")
		os.Exit(1)
	}
	if inputA < 0 || inputB < 0 {
		fmt.Printf("Please provide positive integer values as inputs. Euclidean algorithm doesn't handle negative values well, sadly.")
		os.Exit(1)
	}
	if inputA == 0 && inputB == 0 {
		fmt.Println("at least one input must be different than 0")
		os.Exit(1)
	}

	fmt.Printf("The greatest common divisor between %d and %d is %d\n", inputA, inputB, greatestCommonDivisor(inputA, inputB))
}

func greatestCommonDivisor(x int64, y int64) int64 {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}
