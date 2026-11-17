package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flag.Parse()

	if err := run(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	if flag.NArg() == 0 {
		return errors.New("usage: lifecycle [--] <command> [args...]")
	}
	cmd := exec.Command(flag.Arg(0), flag.Args()[1:]...)
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("start child: %w", err)
	}

	status := make(chan error, 1)
	go func() {
		defer close(status)

		err = cmd.Wait()
		status <- err
	}()

	fmt.Println(cmd.Process.Pid)

	fmt.Printf("Process PID: %d\n", os.Getpid())
	fmt.Printf("Process PPID: %d\n", os.Getppid())

	input := make(chan []byte, 1)

	go func() {
		defer close(input)

		result, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
		if err != nil {
			return
		}
		input <- result
	}()

	select {
	case status := <-status:
		if status != nil {
			return fmt.Errorf("child exit with error: %w", err)
		}
	case <-ctx.Done():
		fmt.Println(" -- Interrupted by User")
		return nil
	case result, ok := <-input:
		if !ok {
			fmt.Println("User no longer chatty")
		} else {
			fmt.Printf("Received: %s", result)
		}
		return nil
	}
	return nil
}
