package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: calc <int>...")
		fmt.Println("usage: calc 2 6 34 -2 55")
		os.Exit(1)
	}

	inputArray, err := getArgArray(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	c := make(chan int)
	go sum(inputArray[:len(inputArray)/3], c)
	go sum(inputArray[len(inputArray)/3:len(inputArray)/3*2], c)
	go sum(inputArray[len(inputArray)/3*2:], c)

	var total int
	for range 3 {
		total += <-c
	}
	fmt.Printf("Sum: %d\n", total)

}

func sum(input []int, c chan int) {
	var total int
	for _, number := range input {
		total += number
	}
	c <- total
}

func getArgArray(stringArray []string) ([]int, error) {
	var array []int
	for _, arg := range stringArray {
		number, err := strconv.Atoi(arg)
		if err != nil {
			return nil, fmt.Errorf("the input %v is not a valid number\n", arg)
		}
		array = append(array, number)
	}
	return array, nil
}
