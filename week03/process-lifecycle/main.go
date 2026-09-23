package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := run(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	fmt.Printf("Process PID: %d\n", os.Getpid())
	fmt.Printf("Process PPID: %d\n", os.Getppid())

	for _, arg := range os.Args[1:] {
		if value, isSet := os.LookupEnv(arg); isSet {
			fmt.Printf("Env var: %s = %s\n", arg, value)
		}
	}

	input := make(chan []byte, 1)

	go func() {
		result, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
		if err != nil {
			return
		}
		input <- result
	}()

	select {
	case <-ctx.Done():
		fmt.Println(" -- Interrupted by User")
		return nil
	case result := <-input:
		fmt.Printf("Received: %s", result)
		return nil
	}
}
