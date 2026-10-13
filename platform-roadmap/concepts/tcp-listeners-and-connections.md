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

- **Deadlines** — `conn.SetReadDeadline(t)` takes an absolute `time.Time`
  (not a duration) and stays set for every later `Read` until changed. For an
  *idle* timeout, push it forward before each read:
  `conn.SetReadDeadline(time.Now().Add(idleTimeout))`. Set once after `Accept`,
  it's a connection lifetime limit instead. `time.Time{}` clears it.
- **Telling read errors apart** — `errors.Is(err, os.ErrDeadlineExceeded)`
  (timeout) vs `errors.Is(err, io.EOF)` (peer closed) vs anything else.
- **Deadlines vs `context`** — deadlines make a blocked `Read`/`Write` return;
  a context does nothing to a goroutine stuck in `conn.Read`. To stop it,
  close the conn or set the deadline to now.

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
- Setting the deadline once after `Accept` and calling it an idle timeout.
- One message for every read error ("client disconnected" for a timeout).
- Returning from the accept loop because one connection's setup failed —
  that shuts down the whole server.
- Forgetting write deadlines: a client that never reads can stall `Write`.

## Seen in
- TCP echo server/client demo, 2026-09
- `week03/tcp-echo-server/server` idle timeout, [2026-10-13](../days/2026-10-13.md)
