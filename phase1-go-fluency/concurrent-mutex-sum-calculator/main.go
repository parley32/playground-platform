package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
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
	var mu sync.Mutex
	var wg sync.WaitGroup
	var result int

	for i := 0; i < len(inputArray); i++ {
		wg.Add(1)
		go func(number int) {
			defer wg.Done()
			mu.Lock()
			result += number
			mu.Unlock()
		}(inputArray[i])
	}
	wg.Wait()
	fmt.Fprintf(output, "Sum: %d\n", result)
	return nil
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
