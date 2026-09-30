# Mutexes, data races and WaitGroups

**Summary:** A mutex makes a section of code run one goroutine at a time so
shared data isn't corrupted; a WaitGroup makes one goroutine wait until others
have finished. They solve different problems.

## Key points
- **Data race** — `result += n` is really read, add, write back. Two goroutines
  can read the same old value, so one addition is lost and the sum changes
  between runs. Behaviour is undefined in Go.
- **Mutex** — `mu.Lock()` / `mu.Unlock()` guard a *section of code*, not the
  variable itself. It only works if every goroutine that touches the data takes
  the lock. Keep the lock held as briefly as possible; do slow work outside it.
- **WaitGroup** — a counter of goroutines still running. `Add(1)` before `go`,
  `Done()` inside (use `defer`), `Wait()` blocks until zero. Without `Wait`,
  `main` exits early and prints an incomplete result.
- **`Wait()` also makes results visible** — reading shared data after `Wait()`
  is safe without a lock.
- **Mutex vs channel** — mutex protects *shared state*; channels *pass data or
  signal* between goroutines (ownership moves, nothing shared). For one shared
  integer a mutex is simpler; for "workers produce, one collector combines" a
  channel fits better.

## Minimal example
```go
var mu sync.Mutex
var wg sync.WaitGroup
var result int

for _, n := range nums {
    wg.Add(1)
    go func(n int) {
        defer wg.Done()
        mu.Lock()
        result += n
        mu.Unlock()
    }(n)
}
wg.Wait()
fmt.Println(result)
```

## Common mistakes
- Calling `wg.Add` inside the goroutine — `Wait` may run first and return early.
- Forgetting `Done()` — `Wait` blocks forever. Use `defer wg.Done()`.
- Some code path touching the data without taking the lock.
- Holding a lock during slow work (network, disk).
- Spawning a goroutine per tiny task: overhead exceeds the work.
- Debug with `go run -race` / `go test -race` — it names the racing lines.

## Not covered yet
`sync.RWMutex`, `sync.Once`, `sync/atomic`.

## Seen in
- `week02/concurrent-mutex-sum-calculator`
- `week02/concurrent-channel-sum-calculator` (channel version)
- [2026-09-30 day log](../days/2026-09-30.md)
