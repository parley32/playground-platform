# Graceful shutdown

**Summary:** Catch the quit signal instead of dying. Stop accepting new
connections, tell the open ones you're shutting down, wake any read that's
blocked waiting for a client instead of waiting out the idle timeout, then wait
for every handler to finish before exiting.

## Key points
- **`signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)`**: the signal
  cancels a context instead of killing the process. Call `stop()` to get default
  signal handling back, so that a second Ctrl+C kills the process.
- **Order:** close the listener (no new connections) → wake every connection →
  `wg.Wait()` → return.
- **Waking a blocked read:** a `select` on `ctx.Done()` in the loop never runs
  while `ReadSlice` is blocked. From *outside* the goroutine, call
  `conn.SetReadDeadline(time.Now())`, and the read returns `os.ErrDeadlineExceeded`.
- **Shutdown or idle timeout?** Both give the same read error. If `ctx.Err() != nil`,
  it's shutdown.
- **Gap:** setting a fresh deadline at the top of each loop can overwrite the
  shutdown deadline, so check `ctx.Err()` right after setting it.
- **`wg.Wait()`** waits for the goroutines started with that WaitGroup
  (`wg.Go`), not for every goroutine.

## Minimal example
```go
stop := context.AfterFunc(ctx, func() { conn.SetReadDeadline(time.Now()) })
defer stop()
// in the read loop, on os.ErrDeadlineExceeded:
if ctx.Err() != nil { /* send notice, return */ }
```

## Common mistakes
- `SetReadDeadline(time.Time{})`: the zero time means **no deadline**, not "now".
- A `go func(){ <-ctx.Done(); ... }()` per connection leaks; use `context.AfterFunc` + `defer stop()`.
- A test whose idle timeout is shorter than its wait passes without the wake-up code.
- A shutdown test proves that connections are woken up, but not that `serve`
  waits for them (`wg.Wait()`). Check that by hand.

Stopping *another* process (TERM → grace → KILL): see [stopping-child-processes](stopping-child-processes.md).

## Seen in
- `phase3-linux-networking/tcp-echo-server/server` (`serve`, `handleConn`, `TestGracefulShutdown`)
- [2026-11-16 day log](../days/2026-11-16.md)
