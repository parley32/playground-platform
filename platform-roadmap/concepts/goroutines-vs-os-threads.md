# Goroutines vs OS threads

**Summary:** Goroutines are cheap because most of their waiting happens in the
Go runtime, not the kernel. A goroutine blocked on a channel is just a struct; a
goroutine blocked in a syscall holds a whole OS thread.

## Key points
- **Linux view:** a "PID" is a thread group ID. Each thread has its own TID
  (the `LWP` column in `ps -L`), and the main thread's TID equals the PID. See
  `/proc/<pid>/task/` and `grep Threads /proc/<pid>/status`.
- **G / M / P:** goroutines (G) run on OS threads (M); an M needs a P (one of
  `GOMAXPROCS` slots) to run Go code. `sysmon` runs without a P.
- **Threads are created on demand and never destroyed.** Idle ones are parked
  on a futex. The count = the most parallel work so far, not a fixed pool.
- **How a goroutine waits** (as seen in `ps -L -o lwp,stat,wchan:20`):

| Waits on | Thread held? | `wchan` |
|---|---|---|
| Channel, mutex, `time.Sleep` | No: a struct in a wait queue | — |
| Network I/O, timers | No: one shared netpoller thread | `ep_poll` |
| Blocking syscall (`wait4`, disk read) | Yes, until it returns | e.g. `do_wait` |
| Idle runtime thread | (parked for reuse) | `futex_do_wait` |

- **In a syscall the kernel blocks the thread itself**, so the runtime can't move
  the goroutine off it. It moves the thread's P to another thread instead, so
  Go code keeps running. That's why the count can exceed `GOMAXPROCS`.
- Sockets are non-blocking under the hood, so `conn.Read` parks the goroutine and
  registers the fd with epoll instead of blocking a thread.

## Common mistakes
- Thinking a thread "watches" blocked goroutines. A send wakes the waiter
  directly: onto a run queue, then onto any free thread.
- Reading the thread count as a fixed default. `GOMAXPROCS=1` lowered it from 8 to 6.
- Assuming lots of blocking syscalls are as cheap as lots of channel waits.
  Each one costs a thread (the default limit is 10,000, then the program crashes).

## Seen in
- `week03/process-lifecycle` (the `cmd.Wait()` thread in `do_wait`)
- [2026-11-23 day log](../days/2026-11-23.md)
