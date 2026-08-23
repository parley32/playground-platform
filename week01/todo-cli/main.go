package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
)

const (
	CMD_ADD    = "add"
	CMD_REMOVE = "remove"
	CMD_LIST   = "list"
	CMD_DONE   = "done"
)

const fileName = "todo.json"

type Todo struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
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
		fmt.Println("Error opening and creating the .json")
		os.Exit(1)
	}
	defer file.Close()

	switch os.Args[1] {
	case CMD_ADD:
		if len(os.Args) != 3 {
			handleInvalidInput()
		}
		if err := addItemToList(file, os.Args[2]); err != nil {
			fmt.Println(err)
		}
	case CMD_REMOVE:
		if len(os.Args) != 3 {
			handleInvalidInput()
		}
		if err := removeItemFromList(file, os.Args[2]); err != nil {
			fmt.Println(err)
		}
	case CMD_LIST:
		if len(os.Args) != 2 {
			handleInvalidInput()
		}
		if err := listItems(file); err != nil {
			fmt.Println(err)
		}
	case CMD_DONE:
		if len(os.Args) != 3 {
			handleInvalidInput()
		}
		if err := checkItem(file, os.Args[2]); err != nil {
			fmt.Println(err)
		}
	}
}

func handleInvalidInput() {
	fmt.Println("usage: go run . <add> <taskName> || <list> || <done | remove> <taskID>")
	os.Exit(1)
}

func addItemToList(file *os.File, name string) (string error) {
	content, err := os.ReadFile(file.Name())
	if err != nil {
		return err
	}

	var todos []Todo
	if err := json.Unmarshal(content, &todos); err != nil {
		return err
	}

	item := Todo{Name: name, Completed: false}

	var highestCurrentID int

	for _, todo := range todos {
		if highestCurrentID, err = strconv.Atoi(todo.ID); err != nil {
			fmt.Printf("the failing id: %v\n", todo.ID)
			return errors.Join(errors.New("Couldnt find highest used id"), err)
		}
	}

	item.ID = strconv.Itoa(highestCurrentID + 1)
	todos = append(todos, item)

	if err := marshalJson(content, &todos); err != nil {
		return err
	}
	fmt.Printf("task number %v: \"%v\" added successfully!\n", item.ID, item.Name)
	return nil
}

func removeItemFromList(file *os.File, id string) error {
	content, err := os.ReadFile(file.Name())
	if err != nil {
		return err
	}

	var todos []Todo
	if err := json.Unmarshal(content, &todos); err != nil {
		return err
	}

	idNumber, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("id not a valid number, list of tasks unchanged")
	}

	if idNumber <= 0 {
		return errors.New("id must be a positive integer, list of tasks unchanged")

	}

	var name string
	taskFound := false
	for i, todo := range todos {
		if todo.ID == id {
			taskFound = true
			name = todo.Name
			todos = append(todos[:i], todos[(i+1):]...)
			break
		}
	}

	if taskFound == false {
		return errors.New("task with that id was not found, list of tasks unchanged")
	}

	if err := marshalJson(content, &todos); err != nil {
		return err
	}
	fmt.Printf("task number %v: \"%v\" removed successfully!\n", id, name)
	return nil

}

func marshalJson(content []byte, todos *[]Todo) error {
	content, err := json.Marshal(todos)
	if err != nil {
		return err
	}

	if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
		return err
	}
	return nil
}

func listItems(file *os.File) error {
	content, err := os.ReadFile(file.Name())
	if err != nil {
		return err
	}

	var todos []Todo
	if err := json.Unmarshal(content, &todos); err != nil {
		return err
	}

	for _, todo := range todos {
		fmt.Printf("%v.", todo.ID)
		if todo.Completed == true {
			fmt.Printf("[X]")
		} else {
			fmt.Printf("[ ]")
		}
		fmt.Printf("%v\n", todo.Name)

	}

	return nil
}

func checkItem(file *os.File, id string) error {
	content, err := os.ReadFile(file.Name())
	if err != nil {
		return err
	}

	var todos []Todo
	if err := json.Unmarshal(content, &todos); err != nil {
		return err
	}

	if rank := slices.IndexFunc(todos, func(T Todo) bool {
		return T.ID == id
	}); rank < 1 {
		return errors.New("No item of that rank was found, list remains unchanged")
	} else {
		todos[rank].Completed = true
		fmt.Printf("task number %v: \"%v\" successfully completed!\n", todos[rank].ID, todos[rank].Name)
	}

	if err := marshalJson(content, &todos); err != nil {
		return err
	}

	return nil
}
