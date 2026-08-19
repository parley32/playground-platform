package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run . <number>")
		fmt.Println("example: go run . 731")
		os.Exit(1)
	}

	number, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		fmt.Println("param is not a valid number")
		os.Exit(1)
	}
	if number <= 1 {
		fmt.Println("param must be a positive number greater than 1")
		os.Exit(1)
	}
	isPrimeNumber(number)

}

func isPrimeNumber(number int64) {
	for i := int64(2); i <= number/i; i++ {
		if number%i == 0 {
			divisor := number / i
			fmt.Printf("%d is not a prime number, it can be written as the product of %d x %d\n", number, i, divisor)
			os.Exit(1)
		}
	}
	fmt.Printf("Congrats, %d is a prime number\n", number)
}
