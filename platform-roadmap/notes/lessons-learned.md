# Lessons learned

Durable gotchas, decisions and "next time do X". Newest at top.
Format: `YYYY-MM-DD — lesson (link to day file or code if useful)`

---

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
