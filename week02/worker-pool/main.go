package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	if err := run(os.Args, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(input []string, output io.Writer) error {
	if len(input) != 3 {
		return fmt.Errorf("Usage: workerPool <workerNumber> <jobsNumber>\nExample: workerPool 3 10")
	}
	workerNumber, err := strconv.Atoi(input[1])
	if err != nil {
		return fmt.Errorf("input %q is not a valid positive number of workers", input[1])
	}
	if workerNumber < 1 {
		return fmt.Errorf("input %q is not a valid positive number of workers", input[1])
	}
	jobsNumber, err := strconv.Atoi(input[2])
	if err != nil {
		return fmt.Errorf("input %q is not a valid positive number of jobs", input[2])
	}
	if jobsNumber < 1 {
		return fmt.Errorf("input %q is not a valid positive number of jobs", input[2])
	}

	if workerNumber > jobsNumber {
		workerNumber = jobsNumber
	}
	results := make(chan int, jobsNumber)
	jobs := make(chan int, jobsNumber)
	var wg sync.WaitGroup

	for i := range workerNumber {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			getReadyForWork(id+1, jobs, results)
		}(i)
	}
	for j := range jobsNumber {
		jobs <- jobsNumber - j
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Fprintf(output, "Result: %v\n", result)
	}

	return nil
}

func getReadyForWork(i int, jobs <-chan int, results chan int) {
	for j := range jobs {
		doJob(i, j, results)
	}
}

func doJob(i int, j int, results chan<- int) {
	time.Sleep(100 * time.Millisecond)
	results <- j
}
