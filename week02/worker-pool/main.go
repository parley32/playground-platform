package main

import (
	"context"
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
	jobsNumber, err := getJobsCount(input[2])
	if err != nil {
		return err
	}
	workerNumber, err := getWorkerCount(input[1], jobsNumber)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	results := make(chan int, jobsNumber)
	jobs := make(chan int, jobsNumber)
	var wg sync.WaitGroup

	for i := range workerNumber {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			getReadyForWork(ctx, jobs, results)
		}(i)
	}
	for j := range jobsNumber {
		jobs <- jobsNumber - j
	}
	close(jobs)

	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel()
		fmt.Fprintln(output, "Job terminated")

	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Fprintf(output, "Result: %v\n", result)
	}

	return nil
}

func getReadyForWork(ctx context.Context, jobs <-chan int, results chan int) {

	for j := range jobs {
		select {
		case <-ctx.Done():
			return
		default:
			doJob(ctx, j, results)
		}
	}
}

func doJob(ctx context.Context, j int, results chan<- int) {
	timer := time.NewTimer(100 * time.Millisecond)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return
	case <-timer.C:
	}

	select {
	case <-ctx.Done():
		return
	case results <- j:
	}
}

func getWorkerCount(input string, jobCount int) (int, error) {
	result, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("input %q is not a valid positive number of workers", input)
	}
	if result < 1 {
		return 0, fmt.Errorf("input %q is not a valid positive number of workers", input)
	}

	if result > jobCount {
		result = jobCount
	}
	return result, nil
}

func getJobsCount(input string) (int, error) {
	result, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("input %q is not a valid positive number of jobs", input)
	}
	if result < 1 {
		return 0, fmt.Errorf("input %q is not a valid positive number of jobs", input)
	}
	return result, nil
}
