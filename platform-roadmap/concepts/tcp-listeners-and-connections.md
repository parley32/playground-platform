# TCP listeners and connections

**Summary:** An IP gets you to a machine, a port to a process on it. A
listener waits on a port; `Accept` turns each incoming client into a
`net.Conn`, a two-way byte stream.

## Key points
- `net.Listen("tcp", ":8080")` → listener. `ln.Accept()` blocks until a
  client connects.
- Handle each connection in its own goroutine so the accept loop keeps going.
- `net.Dial("tcp", "localhost:8080")` is the client side.
- `net.Conn` is an `io.Reader`/`io.Writer`. TCP is a **stream**, not
  messages — you add framing yourself (newline-delimited, length-prefixed…).
- `conn.Close()` releases the fd and tells the peer you're done; a `Read` on
  the other side then returns an error (`io.EOF`).

## Minimal example
```go
ln, _ := net.Listen("tcp", ":8080")
for {
    conn, err := ln.Accept()
    if err != nil { continue }
    go func(c net.Conn) {
        defer c.Close()
        io.Copy(c, c) // echo
    }(conn)
}
```

## Common mistakes
- Not closing connections → fd leak.
- Assuming one `Read` = one message from the sender.
- No read/write deadlines → a silent client holds a goroutine forever
  (`conn.SetDeadline`).

## Seen in
- TCP echo server/client demo, 2026-09
