# Lessons learned

Durable gotchas, decisions and "next time do X". Newest at top.
Format: `YYYY-MM-DD — lesson (link to day file or code if useful)`

---

- 2026-11-24 — Cancel only reaches a goroutine that's waiting inside a `select`
  with a `ctx.Done()` case. A `for range ch` loop blocks outside any select, so
  for a cancellable worker write `for { select { case <-ctx.Done(): return;
  case v, ok := <-ch: if !ok { return } ... } }`. Without `ok`, a closed channel
  hands out zero values forever. ([day](../days/2026-11-24.md))
- 2026-11-23 — A bounded wait is a `select` between the result channel and a
  timer, and after the forced stop you wait on the result again (to get the real
  status and reap the child). Only the sender closes a channel, never the receiver.
- 2026-11-23 — To test a wrapper's signal handling, signal the wrapper's PID with
  `kill`. Ctrl+C goes to the terminal's foreground process group, which may be
  the child and not the wrapper.

- 2026-11-19 — `exec.Cmd` is configuration that `Start()` reads once: set
  `SysProcAttr`, `Stdin`, `Stdout`, `Stderr` and `Env` before it. Setting them
  after has no effect (and is a race). A `nil` `Stdin`/`Stdout`/`Stderr` means
  `/dev/null`, not "inherit". Test with a child that actually reads and prints.
- 2026-11-19 — Sending SIGTERM doesn't guarantee the child exits. A stopped
  child (STAT `T`) ignores it until SIGCONT, and a child can trap it. Send TERM,
  then CONT, and plan for KILL after a grace period, or the supervisor waits
  forever.

- 2026-11-16 — In a test, make the timeout you are *not* testing clearly longer
  than the test's own wait limit. Three times now (2026-10-26, 2026-11-16) a
  short idle timeout made a test pass for the wrong reason. Run a mutation
  check (remove the code under test) to confirm the test can fail.
- 2026-11-16 — A passing mutation check can also show what a test *can't*
  see: removing `wg.Wait()` from `serve` didn't fail `TestGracefulShutdown`,
  because nothing observes the ordering. Note it and check by hand, or add a hook.
- 2026-11-16 — Two Go installs: `export PATH="$PATH:/usr/local/go/bin"` loses to
  the apt `go` in `/usr/bin`. Prepend instead, and check with `which -a go`.
- 2026-09-30 — Two files with `func main()` in one folder work with
  `go run file.go`, but break `go build .` / `go run .` ("main redeclared").
  Give each program its own folder (`server/main.go`, `client/main.go`) when
  you want both to build cleanly.
