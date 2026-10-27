# Testing: table tests, subtests, -race and -cover

**Summary:** Table-driven tests loop one test body over a slice of cases, so a
new scenario is one more row. `-race` finds data races; `-cover` shows which
code the tests execute, not whether they check it.

## Key points
- **Table-driven test** — `tests := []struct{ name string; input ...; want ...; wantErr ... }{...}`
  and a `for` loop running the same body. Adding a case = adding a row.
- **Subtests** — `t.Run(tc.name, func(t *testing.T){...})` names each case, so
  failures read `TestAdd/empty_name`; run one with `go test -run 'TestAdd/empty_name'`.
- **Test the error paths** — rows can expect an error or a specific sentinel
  (`errors.Is`), not just good output.
- **Mocks via interfaces** — pass a fake implementing the interface (see
  [interfaces](interfaces.md)) to test one layer in isolation.
- **`go test -race`** — run it on anything with goroutines; it names the racing lines.
- **`go test -cover`** — percentage of *statements executed* by the tests.
  - See which lines: `go test -coverprofile=c.out ./...` then
    `go tool cover -html=c.out` (or `-func=c.out`).
  - It doesn't tell you whether anything was *asserted*, whether values were
    right, or whether there are races. Use it to find untested code, not as a
    quality score.

- **Test network code with `net.Pipe`** — an in-memory `net.Conn` pair, no real
  port. Closing one end makes reads on the other return `io.EOF`.
- **`net.Pipe` is not real TCP** — `SetReadDeadline` on an already-closed pipe
  returns an error, whereas on real TCP it succeeds and the next read returns
  `io.EOF`. A test that closes the peer immediately can therefore take a
  different code path from run to run (seen: 2 of 12 runs hit the EOF branch).
  For a deterministic "peer closed while we were reading" test, use a loopback
  listener (`net.Listen("tcp", "127.0.0.1:0")` + `net.Dial`).
- **Check a test can fail** — temporarily remove the code under test (or the
  action it reacts to) and confirm the test fails; a test that passes anyway
  is passing for the wrong reason.
- **Assert behaviour, not log output** — logs go to stdout, not the connection,
  so a reader on the client end never sees them. Check what the peer can
  observe (e.g. `io.EOF` after an idle timeout), or inject an `io.Writer` for
  logging if the message itself matters.
- **Deadlines as a test safety net** — `clientConn.SetReadDeadline(...)` makes a
  broken server fail the test quickly instead of hanging it.
- **Avoid `time.Sleep` for synchronisation** — a blocking read already waits;
  a sleep adds flakiness. Make time injectable (a parameter) so tests can use
  milliseconds instead of the production value.
- **A goroutine started by a test should be waited for** (done channel or
  `WaitGroup`), or it can outlive the test.

## Minimal example
```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr error
    }{
        {"valid", "buy milk", nil},
        {"empty", "", ErrInvalidTaskName},
    }
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            _, err := svc.Add(tc.input)
            if !errors.Is(err, tc.wantErr) {
                t.Fatalf("got %v, want %v", err, tc.wantErr)
            }
        })
    }
}
```

## Common mistakes
- Chasing 100% coverage with tests that assert nothing.
- Skipping error cases.
- Shared mutable state between rows.
- Not running `-race` on concurrent code.

## Not covered yet
`t.Parallel`, `t.Helper`, `t.Cleanup`, `httptest`, benchmarks, fuzzing,
golden files.

## Seen in
- `week02/todo-http-api` (`tasks_test.go`, `handler_test.go` with `mockTaskService`)
- `week02/worker-pool`, `week02/concurrent-url-checker`
- [2026-10-06 day log](../days/2026-10-06.md)
