package main

import (
	"fmt"
	"os"
	"strconv"
)

const (
	cmdAdd    = "add"
	cmdRemove = "remove"
	cmdList   = "list"
	cmdDone   = "done"
)

const fileName = "todo.json"

type Todo struct {
	Name      string `json:"name"`
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run . <add> <taskName> || <list> || <done> <taskID> || <remove> <taskId>")
		fmt.Println("example: go run . add \"Add ruleset to Program\"")

		os.Exit(1)
	}

	file, err := os.OpenFile(fileName, os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	switch os.Args[1] {
	case cmdAdd:
		if len(os.Args) != 3 {
			handleInvalidInput()
		}
		success, err := addItemToList(file, os.Args[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(success)
	case cmdRemove:
		if len(os.Args) != 3 {
			handleInvalidInput()
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		success, err := removeItemFromList(file, id)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(success)
	case cmdList:
		if len(os.Args) != 2 {
			handleInvalidInput()
		}
		todos, err := listItems(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, todo := range todos {
			mark := "[ ]"
			if todo.Completed {
				mark = "[X]"
			}
			fmt.Printf("%d.%v%v\n", todo.ID, mark, todo.Name)
		}
	case cmdDone:
		if len(os.Args) != 3 {
			handleInvalidInput()
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		success, err := checkItem(file, id)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(success)
	default:
		handleInvalidInput()

	}
}

func handleInvalidInput() {
	fmt.Println("usage: go run . <add> <taskName> || <list> || <done | remove> <taskID>")
	os.Exit(1)
}
