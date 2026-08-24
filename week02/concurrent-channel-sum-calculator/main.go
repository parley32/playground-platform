package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

const WorkerNumber = 3

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(input []string, output io.Writer) error {
	if len(input) < 1 {
		return fmt.Errorf("Usage: calc <int>...\nExample: calc 2 6 34 -2 55")
	}

	inputArray, err := parseArgs(input)
	if err != nil {
		return err
	}
	c := make(chan int)
	go sum(inputArray[:len(inputArray)/3], c)
	go sum(inputArray[len(inputArray)/3:len(inputArray)/3*2], c)
	go sum(inputArray[len(inputArray)/3*2:], c)

	var total int
	for range WorkerNumber {
		total += <-c
	}
	fmt.Fprintf(output, "Sum: %d\n", total)
	return nil
}

func sum(input []int, c chan int) {
	var total int
	for _, number := range input {
		total += number
	}
	c <- total
}

func parseArgs(stringArray []string) ([]int, error) {
	var array []int
	for _, arg := range stringArray {
		number, err := strconv.Atoi(arg)
		if err != nil {
			return nil, fmt.Errorf("the input %q is not a valid number", arg)
		}
		array = append(array, number)
	}
	return array, nil
}
