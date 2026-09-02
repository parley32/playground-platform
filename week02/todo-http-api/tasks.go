package main

import (
	"fmt"
	"slices"
	"strconv"
	"sync"
)

type TaskService struct {
	mu sync.Mutex
}

func NewTaskService() *TaskService {
	return &TaskService{}
}

func (t *TaskService) ListSingleTask(input string) (Task, error) {
	ID, err := strconv.Atoi(input)
	if err != nil {
		return Task{}, fmt.Errorf("%w: input must be a valid ID, got %q", ErrInvalidID, input)
	}
	if ID < 0 {
		return Task{}, fmt.Errorf("%w:input must be a valid, positive ID, got %q", ErrInvalidID, ID)
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	tasks, err := getStore()
	if err != nil {
		return Task{}, err
	}
	taskRank, err := taskRank(ID, tasks)
	if err != nil {
		return Task{}, err
	}
	return tasks[taskRank], nil
}

func (t *TaskService) List() ([]Task, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return getStore()
}

func (t *TaskService) Add(input string) (Task, error) {
	if len(input) == 0 {
		return Task{}, fmt.Errorf("%w:input must be a valid string, got %q", ErrInvalidTaskName, input)
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	tasks, err := getStore()
	if err != nil {
		return Task{}, err
	}
	newTask := Task{ID: newTaskID(tasks), Text: input, Done: false}
	err = saveStore(append(tasks, newTask))
	if err != nil {
		return Task{}, err
	}
	return newTask, nil
}

func newTaskID(tasks []Task) int {
	var ID int
	for _, task := range tasks {
		if ID <= task.ID {
			ID = task.ID + 1
		}
	}
	return ID
}

func (t *TaskService) Delete(input string) error {
	ID, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("%w:input must be a valid ID, got %q", ErrInvalidID, input)
	}
	if ID < 0 {
		return fmt.Errorf("%w:input must be a valid positive ID, got %q", ErrInvalidID, ID)
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	tasks, err := getStore()
	if err != nil {
		return err
	}
	taskRank, err := taskRank(ID, tasks)
	if err != nil {
		return err
	}
	return saveStore(slices.Delete(tasks, taskRank, taskRank+1))

}

func (t *TaskService) Check(input string) error {
	ID, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("%w:input must be a valid ID, got %q", ErrInvalidID, input)
	}
	if ID < 0 {
		return fmt.Errorf("%w:input must be a valid, positive ID, got %d", ErrInvalidID, ID)
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	tasks, err := getStore()
	if err != nil {
		return err
	}
	taskRank, err := taskRank(ID, tasks)
	if err != nil {
		return err
	}
	if tasks[taskRank].Done == true {
		return fmt.Errorf("%w:task ID %d: %q is already checked", ErrAlreadyDone, ID, tasks[taskRank].Text)
	}
	tasks[taskRank].Done = true
	return saveStore(tasks)
}

func taskRank(ID int, tasks []Task) (int, error) {
	for i, task := range tasks {
		if task.ID == ID {
			return i, nil
		}
	}
	return 0, fmt.Errorf("%w:Could not find task with ID: %q", ErrNotFound, ID)
}
