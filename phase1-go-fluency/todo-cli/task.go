package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
)

func addItemToList(file *os.File, name string) (string, error) {
	todos, err := loadTodos(file)
	if err != nil {
		return "", err
	}

	item := Todo{Name: name, Completed: false}

	item.ID = getHighestCurrentID(todos) + 1
	todos = append(todos, item)

	if err := storeTodos(&todos); err != nil {
		return "", err
	}
	message := fmt.Sprintf("task number %d: %q added successfully!\n", item.ID, item.Name)
	return message, nil
}

func removeItemFromList(file *os.File, id int) (string, error) {
	todos, err := loadTodos(file)
	if err != nil {
		return "", err
	}

	if id <= 0 {
		return "", errors.New("id must be a positive integer, list of tasks unchanged")

	}

	var name string
	taskFound := false
	for i, todo := range todos {
		if todo.ID == id {
			taskFound = true
			name = todo.Name
			todos = slices.Delete(todos, i, i+1)
			break
		}
	}

	if !taskFound {
		return "", errors.New("task with that id was not found, list of tasks unchanged")
	}

	if err := storeTodos(&todos); err != nil {
		return "", err
	}
	message := fmt.Sprintf("task number %v: \"%v\" removed successfully!\n", id, name)
	return message, nil

}

func listItems(file *os.File) ([]Todo, error) {
	return loadTodos(file)
}

func checkItem(file *os.File, id int) (string, error) {
	todos, err := loadTodos(file)
	if err != nil {
		return "", err
	}
	var message string

	if rank := slices.IndexFunc(todos, func(T Todo) bool {
		return T.ID == id
	}); rank < 0 {
		return "", errors.New("no item of that rank was found, list remains unchanged")
	} else if todos[rank].Completed {
		return "", fmt.Errorf("task number %v: \"%v\" is already complete, list remains unchanged", rank+1, todos[rank].Name)
	} else {
		todos[rank].Completed = true
		message = fmt.Sprintf("task %v: \"%v\" successfully completed!\n", todos[rank].ID, todos[rank].Name)
	}

	if err := storeTodos(&todos); err != nil {
		return "", err
	}

	return message, nil
}

func getHighestCurrentID(todos []Todo) int {
	var highestCurrentID int
	for _, todo := range todos {
		if highestCurrentID < todo.ID {
			highestCurrentID = todo.ID
		}
	}
	return highestCurrentID
}
