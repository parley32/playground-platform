package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("usage: go run . <initial sum> <interest percentage> <time in years>")
		fmt.Println("example: go run . 12000 3 10")
		os.Exit(1)
	}

	input := make([]float64, 3)
	for i := 1; i < 4; i++ {
		value, err := strconv.ParseFloat(os.Args[i], 64)
		if err != nil {
			fmt.Printf("Param in position %v is not valid\n", i)
			os.Exit(1)
		} else if value <= 0 {
			fmt.Printf("Param in position %v must be a positive number\n", i)
			os.Exit(1)
		}
		input[i-1] = value
	}
	fmt.Printf("The total value accrued is %.2f\n", calculateInterest(input[0], input[1], input[2]))
}

func calculateInterest(capital float64, interestPercentage float64, years float64) float64 {
	return capital*interestPercentage/100*years + capital
}
