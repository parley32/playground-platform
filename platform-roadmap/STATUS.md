# Status

**Updated:** 2026-11-19
**Phase / week:** Phase 1 complete. Phase 3 (Linux and networking) in progress: first Phase 3 build ("TCP server beyond echo") ticked 2026-11-16. Roadmap re-planned 2026-09-30 toward infrastructure software in Go.
**Project code lives in:** `../week01/`, `../week02/`, `../week03/` (one `go.mod` per week)

## In flight
- `week03/tcp-echo-server`: **done** (steps 1–6: line protocol `PING`/`ECHO`/`SET`/`GET`, length limit, mutex-guarded shared store, idle timeout, graceful shutdown via `signal.NotifyContext` + `context.AfterFunc` + `WaitGroup`; client with separate reader/writer goroutines). Race-tested, checked by hand.
- `week03/process-lifecycle`: plan in `days/2026-11-17.md`. Steps 1–5 done (2026-11-19): child in its own process group, made the terminal's foreground group when stdin is a terminal (`golang.org/x/term`). SIGTERM is sent to the child's group, then SIGCONT, then the wrapper waits. stdio passed through. 128+N exit codes, tested. Orphans and zombies seen in `ps`. Step 6 (threads) left. `week03/go.mod`/`go.sum` (x/term) not committed yet.

## Stuck on
Nothing blocking. Optional follow-ups for the TCP server are in `days/2026-11-16.md` (helper for the duplicated shutdown notice, client EOF wording, Ctrl+D, untested `EOF` branch). Code follow-ups in `week02/todo-http-api` are still listed in `days/2026-09-30.md`.

## Next step
Step 6 of `process-lifecycle`: count the wrapper's threads (`/proc/<pid>/status`, `ps -L`) and explain them, then tick "Processes, threads, scheduling". Small follow-ups in `days/2026-11-19.md` (SIGCONT to the group, diagnostics to stderr, SIGKILL after a grace period). Afterwards: the "worker pool with graceful shutdown" build.

## Background
Hours are not tracked (decision 2026-10-06); pace is judged by ticks and day-file dates. No day files before 2026-09-30, see `days/2026-08-19_to_2026-09-25-summary.md`. Second computer: check that `go version` shows 1.27 (PATH order fix, 2026-11-10).
