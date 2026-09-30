# Interfaces (and dependency injection)

**Summary:** An interface is a set of methods. Any type that has all of them
satisfies it, with no `implements` keyword. Code that depends on the interface
doesn't know or care which concrete type it got.

## Key points
- **Implicit satisfaction** — having the methods is enough. Adding an interface
  later doesn't touch the implementation. Close to TypeScript's structural typing.
- **Depend on the interface, not the concrete type** — `NewTaskHandler` takes a
  `TaskService`, so the handler can't tell the real service from a mock.
- **The interface limits what the caller sees** — a method on the concrete type
  that isn't in the interface is invisible to the handler. To use it, add it to
  the interface, and every implementation (including test mocks) must gain it.
- **Full method set required** — a mock needs all the methods, even ones a test
  never calls, or it won't compile.
- **Benefits beyond testing** — swap implementations without changing callers
  (e.g. `fileStore` today, Postgres later behind the same `Store`).
- **"Accept interfaces, return structs"** — constructors return the concrete
  type; functions that consume dependencies accept an interface.
- **Define the interface where it's used** (the consumer's package), not next
  to the implementation. Optional at small scale.

## Minimal example
```go
type TaskService interface {
    List() ([]Task, error)
}

type TaskHandler struct{ service TaskService }

func NewTaskHandler(s TaskService) *TaskHandler { return &TaskHandler{service: s} }

// In tests:
type mockTaskService struct{ tasks []Task; err error }
func (m mockTaskService) List() ([]Task, error) { return m.tasks, m.err }
```

## Common mistakes
- Making an interface for one implementation with no test or swap reason.
- Huge interfaces: mocks become painful. Prefer small ones.
- Adding a method to the concrete type and expecting the interface holder to see it.

## Not covered yet
Embedding (structs and interfaces), type assertions and type switches, the
nil-interface gotcha.

## Seen in
- `week02/todo-http-api` (`TaskService`, `Store`, `mockTaskService` in `handler_test.go`)
- `week01/shapes`
- [2026-09-30 day log](../days/2026-09-30.md)
