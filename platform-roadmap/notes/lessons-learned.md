# Lessons learned

Durable gotchas, decisions and "next time do X". Newest at top.
Format: `YYYY-MM-DD — lesson (link to day file or code if useful)`

---

- 2026-09-30 — Two files with `func main()` in one folder work with
  `go run file.go`, but break `go build .` / `go run .` ("main redeclared").
  Give each program its own folder (`server/main.go`, `client/main.go`) when
  you want both to build cleanly.
