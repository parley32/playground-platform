# context: cancellation and timeouts

**Summary:** A `context.Context` carries a cancellation signal to concurrent
code. Cancelling closes its `Done()` channel; goroutines stop only because they
choose to check it.

## Key points
- **`ctx, cancel := context.WithCancel(parent)`** creates a cancellable context.
  Always `defer cancel()` — it releases resources even on the normal path.
- **`ctx.Done()` is a channel that gets closed**, not set to true. Receiving
  from a closed channel returns immediately, so one `cancel()` reaches every
  listener at once.
- **Cooperative cancellation** — cancelling doesn't kill goroutines. Each one
  must check `ctx.Done()`, usually in a `select`.
- **Make waits interruptible** — `select` on `ctx.Done()` and a timer/channel
  instead of `time.Sleep`, which can't be interrupted.
- **Guard blocking sends** — `select { case <-ctx.Done(): return; case ch <- v: }`
  so a goroutine isn't stuck forever sending to a channel nobody reads
  (a goroutine leak, which also stops `wg.Wait()` returning).
- **`context.WithTimeout` / `WithDeadline`** cancel automatically after a
  duration; `ctx.Err()` then returns `context.DeadlineExceeded` (vs
  `context.Canceled` for a manual cancel).
- Pass `ctx` as the first parameter of functions that block or do I/O.
- **`stop := context.AfterFunc(ctx, f)`** runs `f` in its own goroutine once
  `ctx` is done, and starts no goroutine before that. `defer stop()` unregisters it
  when the work ends normally. A hand-written `go func(){ <-ctx.Done(); f() }()`
  per connection leaks, because `ctx` may only end at shutdown.
- **`ctx.Err()`** is `nil` until the context is done, and already set when an
  `AfterFunc` callback runs. Use it to find out *why* an operation failed, for
  example a read deadline caused by shutdown rather than an idle timeout.

## Minimal example
```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

select {
case <-ctx.Done():
    return
case <-time.After(100 * time.Millisecond):
    // do the work
}
```

## Common mistakes
- Forgetting `defer cancel()`.
- A goroutine that never checks `ctx.Done()` — cancel has no effect on it.
- `time.Sleep` inside a cancellable worker.
- Unguarded send or receive on a channel that may never be serviced.
- Hand-rolling a timeout with `sleep` + `cancel()` instead of `WithTimeout`.
- A `select` on `ctx.Done()` can't interrupt a call that's already blocked
  (e.g. `conn.Read`). Wake it from outside (`AfterFunc` + `SetReadDeadline(time.Now())`).

## Not covered yet
`WithTimeout`/`WithDeadline` in practice, passing `ctx` into `net/http`
requests, `context.WithValue`.

## Seen in
- `week02/worker-pool` (`WithCancel`, `doJob` select)
- `week02/concurrent-url-checker`
- [2026-09-30 day log](../days/2026-09-30.md)
- `week03/tcp-echo-server/server` (`AfterFunc`, `ctx.Err()`), see [graceful shutdown](graceful-shutdown.md), [2026-11-16](../days/2026-11-16.md)
