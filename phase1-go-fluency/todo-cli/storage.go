package main

import (
	"encoding/json"
	"os"
)

func storeTodos(todos *[]Todo) error {
	content, err := json.Marshal(todos)
	if err != nil {
		return err
	}

	if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
		return err
	}
	return nil
}

func loadTodos(file *os.File) ([]Todo, error) {
	content, err := os.ReadFile(file.Name())
	if err != nil {
		return nil, err
	}

	if len(content) == 0 {
		return []Todo{}, nil
	}

	var todos []Todo
	if err := json.Unmarshal(content, &todos); err != nil {
		return nil, err
	}

	return todos, nil
}
