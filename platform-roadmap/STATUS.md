# Status

**Updated:** 2026-10-29
**Phase / week:** Phase 1 complete (all concept and build boxes ticked, understanding pass done 2026-10-06), with an early start on Phase 3 (Linux and networking) material. Roadmap re-planned 2026-09-30 toward infrastructure software in Go.
**Project code lives in:** `../week01/`, `../week02/`, `../week03/` (one `go.mod` per week)

## In flight
- Phase 3 (Linux and networking): `week03/process-lifecycle` and `week03/tcp-echo-server` (idle timeout, line protocol with `PING`/`ECHO`/`SET`/`GET`, shared mutex-guarded store; steps 1–4 done and race-tested; line limit and graceful shutdown left)

## Stuck on
Not stuck, but a comprehension gap: nothing blocking. All Phase 1 concepts are explained. Open code follow-ups in `week02/todo-http-api` are listed in `days/2026-09-30.md`. `-race` is already in use with test suites. Still unsure about `-cover` and modules/package layout.

## Next step
Phase 3 TCP server step 5: reject lines over a maximum length (`days/2026-10-29.md`), then step 6 (graceful shutdown). Optional: loopback-TCP test for the EOF branch. Afterwards: explain what you did without notes; the agent writes the concept entry from your explanation.

## Background
Hours are not tracked (decision 2026-10-06); pace is judged by ticks and day-file dates. No day files before 2026-09-30 — see `days/2026-08-19_to_2026-09-25-summary.md`.
