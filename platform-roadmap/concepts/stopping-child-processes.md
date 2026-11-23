# Stopping a child process: TERM → grace period → KILL

**Summary:** A supervisor asks politely first (SIGTERM), gives the child a fixed
time to clean up, and only then forces it (SIGKILL). Then it still waits, to get
the exit status and reap the child. Same idea as `docker stop`. (For a server
shutting *itself* down, see [graceful-shutdown](graceful-shutdown.md).)

## Key points
- **Signal the group:** `syscall.Kill(-pid, sig)` reaches grandchildren too
  (needs `Setpgid: true`). `ESRCH` = already gone, not an error.
- **TERM, then CONT:** a stopped child (STAT `T`) queues TERM until it continues.
- **The grace period is a `select`** between "the child exited" and "time's up".
  A bare `<-status` waits forever if the child traps TERM.
- **After KILL, wait again:** the real status (killed by 9 → exit 137) arrives on
  the same channel, and waiting reaps the zombie.
- SIGKILL can't be caught, so the child's `defer`s and handlers don't run.

## Minimal example
```go
syscall.Kill(-pid, syscall.SIGTERM)
syscall.Kill(-pid, syscall.SIGCONT)
select {
case err := <-status:
    return err
case <-time.After(grace):
    syscall.Kill(-pid, syscall.SIGKILL)
    return <-status // 137
}
```

## Common mistakes
- Putting KILL after the wait, so it only runs once the child has already exited.
- `close(status)` from the receiver: only the sender (the `Wait` goroutine) closes.
- Returning `nil` after KILL, so the exit is 0 instead of 137.
- Testing with Ctrl+C when the child has the terminal's foreground: only the child
  gets SIGINT. Use `kill -TERM <wrapper pid>`.

## Seen in
- `week03/process-lifecycle/main.go` (`run`, `7f071a6`)
- [2026-11-19](../days/2026-11-19.md), [2026-11-23](../days/2026-11-23.md)
