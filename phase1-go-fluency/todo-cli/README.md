# Todo CLI

A small command-line todo list written in Go. Tasks are stored locally in a
`todo.json` file in the directory where the program is run.

## Requirements

- Go 1.25 or later

## Run

From this directory, run commands with:

```sh
go run . <command> [argument]
```

### Add a task

```sh
go run . add "Buy groceries"
```

### List tasks

```sh
go run . list
```

Tasks are displayed with their ID and completion status:

```text
1.[ ] Buy groceries
2.[X] Submit assignment
```

### Mark a task complete

```sh
go run . done 1
```

### Remove a task

```sh
go run . remove 1
```

Task IDs are positive integers. Removing a task does not reuse its ID; new
tasks receive an ID one greater than the highest existing ID.

## Build

Create an executable with:

```sh
go build -o todo .
./todo add "Read Go documentation"
```

## Test

Run the package tests and static checks with:

```sh
go test
go vet .
```

The tests use temporary directories, so they do not modify the project's
`todo.json` file.

## Data format

Tasks are stored as a JSON array:

```json
[
	{
		"name": "Buy groceries",
		"id": 1,
		"completed": false
	}
]
```
