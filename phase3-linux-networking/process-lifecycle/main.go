package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/term"
)

var errUsage = errors.New("usage: lifecycle [--] <command> [args...]")

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flag.Parse()

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		statusCode := exitCode(err)
		os.Exit(statusCode)
	}
}

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	if flag.NArg() == 0 {
		return errUsage
	}
	cmd := exec.Command(flag.Arg(0), flag.Args()[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Ctty: 0} // verbose 'cause fd 0 is required for Foreground trick to work

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if term.IsTerminal(int(os.Stdin.Fd())) {
		cmd.SysProcAttr.Foreground = true
	}

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("start child: %w", err)
	}

	status := make(chan error, 1)
	go func() {
		defer close(status)

		err := cmd.Wait()
		status <- err
	}()

	fmt.Fprintln(os.Stderr, cmd.Process.Pid)

	fmt.Fprintf(os.Stderr, "Process PID: %d\n", os.Getpid())
	fmt.Fprintf(os.Stderr, "Process PPID: %d\n", os.Getppid())

	select {
	case err := <-status:
		if err != nil {
			return fmt.Errorf("child exit with error: %w", err)
		}
	case <-ctx.Done():
		cmdErr := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
		if cmdErr != nil && !errors.Is(cmdErr, syscall.ESRCH) {
			return fmt.Errorf("failed to signal child process: %w", cmdErr)
		}
		fmt.Fprintln(os.Stderr, " -- received signal, stopping child")

		cmdErr = syscall.Kill(-cmd.Process.Pid, syscall.SIGCONT) // negative PID means the whole group, which is nice if the child has more children
		if cmdErr != nil && !errors.Is(cmdErr, syscall.ESRCH) {
			return fmt.Errorf("failed to signal child process: %w", cmdErr)
		}

		select {
		case <-time.After(3 * time.Second):
			cmdErr := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			if cmdErr != nil && !errors.Is(cmdErr, syscall.ESRCH) {
				return fmt.Errorf("failed to kill child process: %w", cmdErr)
			}
			if statusErr := <-status; statusErr != nil {
				return fmt.Errorf("child exit with error: %w", statusErr)
			}
		case statusErr := <-status:
			if statusErr != nil {
				return fmt.Errorf("child exit with error: %w", statusErr)
			}
		}
	}
	return nil
}

func exitCode(err error) int {
	if exitErr, ok := errors.AsType[*exec.ExitError](err); ok { // same as declaring a var of type *exec.ExitError, and then checking errors.As on err against that var
		ws, ok := exitErr.Sys().(syscall.WaitStatus)
		if ok && ws.Signaled() {
			return 128 + int(ws.Signal())
		}
		return exitErr.ExitCode()
	}
	if errors.Is(err, exec.ErrNotFound) {
		return 127 // keep exit status code 1 for parent process failure
	}
	if errors.Is(err, errUsage) {
		return 2
	}
	if err == nil {
		return 0
	}
	return 1
}
