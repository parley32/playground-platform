# Status

**Updated:** 2026-10-26
**Phase / week:** Phase 1 complete (all concept and build boxes ticked, understanding pass done 2026-10-06), with an early start on Phase 3 (Linux and networking) material. Roadmap re-planned 2026-09-30 toward infrastructure software in Go.
**Project code lives in:** `../week01/`, `../week02/`, `../week03/` (one `go.mod` per week)

## In flight
- Phase 3 (Linux and networking): `week03/process-lifecycle` and `week03/tcp-echo-server` (echo works, basic handler tests; no framing, concurrency or timeouts yet)

## Stuck on
Not stuck, but a comprehension gap: nothing blocking. All Phase 1 concepts are explained. Open code follow-ups in `week02/todo-http-api` are listed in `days/2026-09-30.md`. `-race` is already in use with test suites. Still unsure about `-cover` and modules/package layout.

## Next step
Phase 3 TCP server beyond echo: see the plan in `days/2026-10-06.md` (idle timeout, line protocol, SET/GET with a mutex, line-length limit, graceful shutdown). Steps 1–2 done; step 3 (table test) passes. Remaining for step 3: tests for the EOF and idle-timeout paths of `handleConn` (make the timeout injectable; coverage of `handleConn` is 60%). Then step 4: `SET`/`GET` with a mutex. Afterwards: explain what you did without notes; the agent writes the concept entry from your explanation.
## Background
Hours are not tracked (decision 2026-10-06); pace is judged by ticks and day-file dates. No day files before 2026-09-30 — see `days/2026-08-19_to_2026-09-25-summary.md`.
