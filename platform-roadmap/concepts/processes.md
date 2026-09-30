# Processes, PIDs and parent/child

**Summary:** A program is a file on disk; a process is a running instance of
it, with its own memory, PID, environment, open files and exit code.

## Key points
- **PID** — unique number the kernel gives each process. `os.Getpid()`.
- **Parent/child** — every process is started by another (PID 1 is the root).
  Children inherit a copy of the parent's environment variables and open
  file descriptors.
- **Environment variables** — key/value strings copied from the parent at
  start. `os.Getenv`, `os.Environ`; set `cmd.Env` to control a child's.
- **Signals** — async notifications. `SIGINT` (Ctrl+C) and `SIGTERM` can be
  caught; `SIGKILL` cannot. Catch with `os/signal.Notify`.
- **Exit codes** — `0` = success, non-zero = failure. `os.Exit(n)`; returning
  from `main` = `0`. Shell shows it with `echo $?`.

## Minimal example
```go
cmd := exec.Command("ls", "-l")
cmd.Stdout = os.Stdout
if err := cmd.Run(); err != nil {
    var exitErr *exec.ExitError
    if errors.As(err, &exitErr) {
        fmt.Println("child exited with", exitErr.ExitCode())
    }
}
```

## Common mistakes
- `os.Exit` skips deferred functions — cleanup won't run.
- Forgetting a child inherits your env and fds (secrets, open sockets).

## Seen in
- Intro conversation, 2026-09
