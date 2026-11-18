# Status

**Updated:** 2026-11-18
**Phase / week:** Phase 1 complete. Phase 3 (Linux and networking) in progress: first Phase 3 build ("TCP server beyond echo") ticked 2026-11-16. Roadmap re-planned 2026-09-30 toward infrastructure software in Go.
**Project code lives in:** `../week01/`, `../week02/`, `../week03/` (one `go.mod` per week)

## In flight
- `week03/tcp-echo-server`: **done** (steps 1–6: line protocol `PING`/`ECHO`/`SET`/`GET`, length limit, mutex-guarded shared store, idle timeout, graceful shutdown via `signal.NotifyContext` + `context.AfterFunc` + `WaitGroup`; client with separate reader/writer goroutines). Race-tested, checked by hand.
- `week03/process-lifecycle`: plan in `days/2026-11-17.md` (6 steps). Steps 1–3 done (child from the command line, PID/PPID checked with `ps`/`pstree`, child's exit code passed through via a tested `exitCode`; `echo $?` shows `3`).

## Stuck on
Nothing blocking. Optional follow-ups for the TCP server are in `days/2026-11-16.md` (helper for the duplicated shutdown notice, client EOF wording, Ctrl+D, untested `EOF` branch). Code follow-ups in `week02/todo-http-api` are still listed in `days/2026-09-30.md`.

## Next step
Step 4 of `process-lifecycle`: forward SIGTERM to the child and `Wait()` before exiting; first check whether Ctrl+C already reaches the child (process groups, `Setpgid`). Also fix: stdin EOF currently wins the `select` and the program exits 0 instead of the child's code. Then steps 5–6. Afterwards: the "worker pool with graceful shutdown" build.

## Background
Hours are not tracked (decision 2026-10-06); pace is judged by ticks and day-file dates. No day files before 2026-09-30, see `days/2026-08-19_to_2026-09-25-summary.md`. Second computer: check that `go version` shows 1.27 (PATH order fix, 2026-11-10).
