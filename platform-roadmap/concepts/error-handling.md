# Error handling: sentinels, wrapping, Is and As

**Summary:** Errors are values. Wrap them with `%w` to add context while
keeping the original inside a chain; use `errors.Is` to ask "is it this
error?" and `errors.As` to ask "is it this kind of error, and give me its
details".

## Key points
- **Sentinel error** — one shared value, `var ErrNotFound = errors.New("...")`,
  that callers compare against.
- **Wrapping** — `fmt.Errorf("context: %w", err)` adds text and keeps `err`
  in a chain. `%v` gives the same text but loses the identity, so `Is`/`As`
  won't find it.
- **`errors.Is(err, target)`** — walks the chain; true if any error in it
  equals `target`. Use for sentinels.
- **`errors.As(err, &target)`** — walks the chain looking for an error whose
  *type* matches `target`'s type; if found, assigns it into `target` and returns
  true. Use to read fields/methods from a specific error type.
- **Map errors at the boundary** — internal code returns sentinels; the HTTP
  layer decides status codes (`SendError`: `ErrNotFound` → 404, etc.).
- **Message convention** — `"%w: detail"` (colon *and space*).

## Minimal example
```go
var ErrNotFound = errors.New("task not found")

func find(id int) error {
    return fmt.Errorf("%w: task ID %d", ErrNotFound, id)
}

err := find(3)
if errors.Is(err, ErrNotFound) { /* 404 */ }

var exitErr *exec.ExitError
if errors.As(err, &exitErr) { fmt.Println(exitErr.ExitCode()) }
```

## Common mistakes
- Using `%v` instead of `%w`, then wondering why `errors.Is` fails.
- `"%w:msg"` with no space after the colon.
- Comparing errors with `==` after they've been wrapped.
- Reusing a sentinel for an unrelated problem (e.g. a malformed JSON body
  reported as `ErrInvalidTaskName`).
- Sending `err.Error()` to API clients for unexpected errors — can leak
  internals. Log the detail, return a generic message.
- Defining a sentinel and never using it (`ErrInternal`).

## Not covered yet
Custom error types (a struct implementing `error`) used with `errors.As`;
`errors.Join`; `panic`/`recover`.

## Seen in
- `week02/todo-http-api` (`SendError`, sentinels in `handler.go`, wrapping in `tasks.go`)
- `week01/config-reader`
- [2026-09-30 day log](../days/2026-09-30.md)
