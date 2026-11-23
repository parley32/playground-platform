# Status

**Updated:** 2026-11-23
**Phase / week:** Phase 1 complete. Phase 3 (Linux and networking) in progress: first Phase 3 build ("TCP server beyond echo") ticked 2026-11-16. Roadmap re-planned 2026-09-30 toward infrastructure software in Go.
**Project code lives in:** `../phase1-go-fluency/`, `../phase3-linux-networking/` (one folder and one `go.mod` per phase; renamed from `week01`–`week03` on 2026-11-23)

## In flight
- `phase3-linux-networking/tcp-echo-server`: **done** (steps 1–6: line protocol `PING`/`ECHO`/`SET`/`GET`, length limit, mutex-guarded shared store, idle timeout, graceful shutdown via `signal.NotifyContext` + `context.AfterFunc` + `WaitGroup`; client with separate reader/writer goroutines). Race-tested, checked by hand.
- `phase3-linux-networking/process-lifecycle`: **done** (2026-11-23). Child in its own process group with the terminal foreground, stdio passed through, TERM + CONT to the group, SIGKILL after a 3 s grace period (`7f071a6`), 128+N exit codes (143/137 tested). Orphans, zombies and the wrapper's threads (`futex`/`do_wait`/`ep_poll`, GOMAXPROCS) observed. Optional: pass on the received signal; a test for the grace period; Ctrl+C can't reach the wrapper while the child has the foreground (design question in `days/2026-11-23.md`).

## Stuck on
Nothing blocking. Optional follow-ups for the TCP server are in `days/2026-11-16.md` (helper for the duplicated shutdown notice, client EOF wording, Ctrl+D, untested `EOF` branch). Code follow-ups in `phase1-go-fluency/todo-http-api` are still listed in `days/2026-09-30.md`.

## Next step
"Worker pool with graceful shutdown" build, planned in `days/2026-11-24.md`. It starts from a copy of `phase1-go-fluency/worker-pool`, whose `go test -race` currently fails (a worker blocked on an empty `jobs` channel ignores cancel; data race at `main.go:55`). First action: `cp -r ../phase1-go-fluency/worker-pool .` from `phase3-linux-networking/`, then `go test -race ./worker-pool`.

## Background
Hours are not tracked (decision 2026-10-06); pace is judged by ticks and day-file dates. No day files before 2026-09-30, see `days/2026-08-19_to_2026-09-25-summary.md`. Second computer: check that `go version` shows 1.27 (PATH order fix, 2026-11-10).
