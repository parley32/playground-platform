# Status

**Updated:** 2026-09-30
**Phase / week:** Phase 1 wrapping up (understanding pass pending), with an early start on Phase 3 (Linux and networking) material. Roadmap re-planned 2026-09-30 toward infrastructure software in Go.
**Project code lives in:** `../week01/`, `../week02/`, `../week03/` (one `go.mod` per week)

## In flight
- Phase 3 (Linux and networking): `week03/process-lifecycle` and `week03/tcp-echo-server` (echo works, basic handler tests; no framing, concurrency or timeouts yet)

## Stuck on
Not stuck, but a comprehension gap: interfaces, mutexes, `context`, error handling and channels are now understood; still to check: modules/package layout, testing (`-cover`, table tests). `-race` is already in use with test suites. Still unsure about `-cover` and modules/package layout.

## Next step
Understanding pass on modules and package layout via `week02/todo-http-api` and the `go.mod` files (module vs package, `package main`, exported names, splitting into packages): explain it back without notes; the agent writes the concept entry from your explanation.
## Background
Weekly hours are well under the roadmap's 34 h baseline, slowly rising. No day files before 2026-09-30 — see `days/2026-08-19_to_2026-09-25-summary.md`.
