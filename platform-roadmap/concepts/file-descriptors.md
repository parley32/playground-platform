# File descriptors and stdin/stdout/stderr

**Summary:** A file descriptor is a small integer indexing a per-process table
of open "files" — real files, pipes, sockets, terminals.

## Key points
- Every process starts with **0 = stdin**, **1 = stdout**, **2 = stderr**.
  New opens get the lowest free number.
- stderr is separate so `prog > out.log` still shows errors on screen.
- Children inherit copies of the parent's fds — that's how `a | b` works
  (a's fd 1 is wired to b's fd 0).
- Fds are finite per process (`ulimit -n`). Leaking them (not closing files
  or connections) eventually breaks the program.
- In Go: `os.Stdin/Stdout/Stderr` are `*os.File`; `f.Fd()` gives the raw
  number. Files and network connections both implement `io.Reader/Writer`.

## Minimal example
```go
f, err := os.Open("data.txt")
if err != nil {
    fmt.Fprintln(os.Stderr, "open failed:", err)
    os.Exit(1)
}
defer f.Close()
io.Copy(os.Stdout, f)
```

## Common mistakes
- Forgetting `defer f.Close()` → fd leak.
- Printing errors to stdout, so they end up mixed into piped output.

## Seen in
- Intro conversation, 2026-09
