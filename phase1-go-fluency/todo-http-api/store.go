package main

import (
	"encoding/json"
	"os"
)

const fileName = "todo.json"

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type Store interface {
	getStore() ([]Task, error)
	saveStore([]Task) error
}

type fileStore struct{}

func (fileStore) getStore() ([]Task, error) {
	return getStore()
}

func (fileStore) saveStore(tasks []Task) error {
	return saveStore(tasks)
}

func getStore() ([]Task, error) {
	content, err := os.ReadFile(fileName)
	if os.IsNotExist(err) {
		return []Task{}, nil
	}
	if err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return []Task{}, nil
	}

	var tasks []Task
	if err := json.Unmarshal(content, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func saveStore(content []Task) error {
	tasks, err := json.Marshal(content)
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, tasks, 0644)
}
