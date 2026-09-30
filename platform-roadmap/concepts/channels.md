# Channels: buffered, unbuffered and closing

**Summary:** A channel passes values between goroutines. Unbuffered channels
synchronise sender and receiver; buffered ones decouple them up to a capacity.
Channels never close themselves — `close(ch)` is an explicit "no more values"
signal to receivers.

## Key points
- **Unbuffered** (`make(chan int)`) — a send blocks until a receiver takes the
  value. Sender and receiver meet.
- **Buffered** (`make(chan int, n)`) — sends don't block until `n` values are
  waiting; then the next send blocks. Receives block when empty.
- **Closing is explicit and is a signal, not cleanup** — the GC handles
  channels. Close so receivers know no more values are coming (`for v := range ch`
  ends when the channel is closed and drained).
- **Who closes** — the sending side, once *all* senders are done. With several
  senders: `go func() { wg.Wait(); close(ch) }()`.
- **Rules** — sending on a closed channel panics; closing twice panics;
  receiving from a closed channel returns the zero value immediately
  (`v, ok := <-ch` gives `ok == false`).
- **Goroutine leak** — a goroutine blocked forever on a send/receive nobody
  will match. Doesn't panic unless *every* goroutine is stuck (then: "all
  goroutines are asleep - deadlock!").
- **Counting instead of closing** — `for range n { total += <-c }` works when
  you know exactly how many values will arrive.
- **Don't `wg.Wait()` in the receiver before ranging** with an unbuffered
  channel — senders wait on the receiver, the receiver waits on the senders.

## Minimal example
```go
c := make(chan int)
var wg sync.WaitGroup
for _, part := range parts {
    wg.Add(1)
    go func() {
        defer wg.Done()
        c <- sum(part)
    }()
}
go func() {
    wg.Wait()
    close(c)
}()

total := 0
for v := range c {
    total += v
}
```

## Common mistakes
- Believing buffered channels close automatically.
- Receiver (not sender) closing the channel, or closing while senders remain.
- `wg.Wait()` in the same goroutine that must receive.
- Receiving fewer values than were sent on an unbuffered channel (leak).

## Not covered yet
`select` with multiple channels (partly seen with `ctx.Done()`), nil channels,
directional channel types in depth (`<-chan`, `chan<-`), fan-in/fan-out patterns.

## Seen in
- `week02/concurrent-channel-sum-calculator` (unbuffered, count receives)
- `week02/worker-pool` (buffered jobs/results, `close(jobs)`, `close(results)` after `wg.Wait()`)
- [2026-09-30 day log](../days/2026-09-30.md)
