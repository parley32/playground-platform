# Modules and packages

**Summary:** A package is one directory of `.go` files compiled together; a
module is a set of packages versioned together, defined by `go.mod`. Capitalised
names are visible to other packages, lowercase ones aren't.

## Key points
- **Package** — all files in one directory share one `package name`. Only
  `package main` builds an executable; any other name is a library. You can't
  import `main`.
- **Module** — a tree of packages with a `go.mod` at the root. `go.mod` holds the
  module **name** (the import-path prefix), the Go version, and any outside
  dependencies.
- **Import path** = module name + folder path, e.g. `"phase1-go-fluency/todo-http-api/handler"`.
- **Exported vs unexported** — a capital first letter makes *any* identifier
  (function, type, constant, variable, struct field, method) visible outside its
  package. Lowercase is package-private.
- **One program per folder** — two `func main()` in one folder is one package
  declaring `main` twice. Give each program its own folder.
- **No import cycles** — if A imports B, B can't import A. When splitting a
  `package main` into packages, shared things (sentinel errors, domain types)
  go in a package both sides import, not in the HTTP layer.
- **Common layout** — `cmd/<app>/main.go` for entry points, other packages
  beside it, `internal/` for code only this module may import.

## Minimal example
```
phase1-go-fluency/     module phase1-go-fluency   (go.mod)
  todo-http-api/        package main
    handler/            package handler  -> import "phase1-go-fluency/todo-http-api/handler"
```
```go
package handler

func NewTaskHandler(s TaskService) *TaskHandler { ... } // visible outside
func parseID(s string) (int, error)            { ... } // package-private
```

## Common mistakes
- Thinking the `go` line in `go.mod` is what defines a module (it's a detail;
  the module name is the point).
- Moving a file to a new folder but leaving `package main` in it.
- Expecting lowercase names to be reachable from another package.
- Putting shared errors/types in the layer that depends on the others, causing
  an import cycle (e.g. sentinels in `handler.go` but also used by `tasks.go`).

## Not covered yet
Adding dependencies (`go get`, `go.sum`), `internal/` in practice, workspaces
(`go.work`), versioning/semver.

## Seen in
- `phase1-go-fluency/go.mod`, `phase3-linux-networking/go.mod` (one module per phase; `week01` + `week02` were merged into one module on 2026-11-23, which only meant deleting a `go.mod` because nothing imported by module name)
- `phase1-go-fluency/todo-http-api` (all `package main`)
- `phase3-linux-networking/tcp-echo-server/server` and `client`
- [2026-10-06 day log](../days/2026-10-06.md)
