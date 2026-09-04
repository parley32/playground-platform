package main

import (
	"errors"
	"reflect"
	"testing"
)

type mockStore struct {
	tasks []Task
}

func (m *mockStore) getStore() ([]Task, error) {
	return m.tasks, nil
}

func (m *mockStore) saveStore(tasks []Task) error {
	m.tasks = tasks
	return nil
}

func TestListSingleTask(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		list    []Task
		want    Task
		wantErr error
	}{
		{
			name:  "success list item",
			input: "1",
			list:  []Task{{ID: 1, Text: "Item 1", Done: false}},
			want:  Task{ID: 1, Text: "Item 1", Done: false},
		},
		{
			name:    "invalid input type",
			input:   "blabla",
			list:    []Task{{ID: 1, Text: "Item 1", Done: false}},
			want:    Task{},
			wantErr: ErrInvalidID,
		},
		{
			name:    "no task with matching input",
			input:   "2",
			list:    []Task{{ID: 1, Text: "Item 1", Done: false}},
			want:    Task{},
			wantErr: ErrNotFound,
		},
		{
			name:    "no existing list",
			input:   "1",
			list:    []Task{},
			want:    Task{},
			wantErr: ErrNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := NewTaskService(&mockStore{tasks: test.list})

			got, err := f.ListSingleTask(test.input)
			if err != nil {
				if !errors.Is(err, test.wantErr) {
					t.Fatalf("Got error %q, expected %q", err, test.wantErr)
				}
			}
			if got != test.want {
				t.Fatalf("Got %v, expected %v", got, test.want)
			}
		})
	}
}

func TestList(t *testing.T) {
	tests := []struct {
		name    string
		list    []Task
		want    []Task
		wantErr error
	}{
		{
			name: "success list item",
			list: []Task{{ID: 1, Text: "Item 1", Done: false}},
			want: []Task{{ID: 1, Text: "Item 1", Done: false}},
		},
		{
			name:    "no existing list",
			list:    []Task{},
			want:    []Task{},
			wantErr: ErrNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := NewTaskService(&mockStore{tasks: test.list})

			got, err := f.List()
			if err != nil {
				if !errors.Is(err, test.wantErr) {
					t.Fatalf("Got error %q, expected %q", err, test.wantErr)
				}
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("Got %v, expected %v", got, test.want)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		list    []Task
		want    Task
		wantErr error
	}{
		{
			name:  "success add item to list",
			input: "Item todo",
			list:  []Task{{ID: 0, Text: "Item 1", Done: false}},
			want:  Task{ID: 1, Text: "Item todo", Done: false},
		},
		{
			name:  "success add item to empty list",
			input: "Item todo",
			list:  []Task{},
			want:  Task{ID: 0, Text: "Item todo", Done: false},
		},
		{
			name:  "success add item to list with missing IDs",
			input: "Item todo",
			list:  []Task{{ID: 10, Text: "Item 1", Done: false}},
			want:  Task{ID: 11, Text: "Item todo", Done: false},
		},
		{
			name:    "empty input",
			list:    []Task{{ID: 0, Text: "Item 1", Done: false}},
			want:    Task{},
			wantErr: ErrInvalidTaskName,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := NewTaskService(&mockStore{tasks: test.list})

			got, err := f.Add(test.input)
			if err != nil {
				if !errors.Is(err, test.wantErr) {
					t.Fatalf("Got error %q, expected %q", err, test.wantErr)
				}
			}
			if got != test.want {
				t.Fatalf("Got %v, expected %v", got, test.want)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		list    []Task
		want    []Task
		wantErr error
	}{
		{
			name:  "success delete item from list",
			input: "1",
			list:  []Task{{ID: 0, Text: "Item 1", Done: false}, {ID: 1, Text: "Item todo", Done: false}},
			want:  []Task{{ID: 0, Text: "Item 1", Done: false}},
		},
		{
			name:  "success delete last item from list",
			input: "1",
			list:  []Task{{ID: 1, Text: "Item todo", Done: false}},
			want:  []Task{},
		},
		{
			name:    "failure no task with that ID",
			input:   "1",
			list:    []Task{{ID: 10, Text: "Item 1", Done: false}},
			want:    []Task{{ID: 10, Text: "Item 1", Done: false}},
			wantErr: ErrNotFound,
		},
		{
			name:    "failure empty input",
			input:   "",
			list:    []Task{{ID: 0, Text: "Item 1", Done: false}},
			want:    []Task{{ID: 0, Text: "Item 1", Done: false}},
			wantErr: ErrInvalidID,
		},
		{
			name:    "failure invalid input",
			input:   "blabla",
			list:    []Task{{ID: 0, Text: "Item 1", Done: false}},
			want:    []Task{{ID: 0, Text: "Item 1", Done: false}},
			wantErr: ErrInvalidID,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := &mockStore{tasks: test.list}
			f := NewTaskService(store)

			err := f.Delete(test.input)
			if err != nil && !errors.Is(err, test.wantErr) {
				t.Fatalf("Got error %q, expected %q", err, test.wantErr)
			}
			if !reflect.DeepEqual(store.tasks, test.want) {
				t.Fatalf("Got %v, expected %v", store.tasks, test.want)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		list    []Task
		want    []Task
		wantErr error
	}{
		{
			name:  "success item completion",
			input: "1",
			list:  []Task{{ID: 1, Text: "Item to check", Done: false}},
			want:  []Task{{ID: 1, Text: "Item to check", Done: true}},
		},
		{
			name:    "failed no item with specified input",
			input:   "2",
			list:    []Task{{ID: 1, Text: "Item to check", Done: false}},
			want:    []Task{{ID: 1, Text: "Item to check", Done: false}},
			wantErr: ErrNotFound,
		},
		{
			name:    "failed, invalid input",
			input:   "bla",
			list:    []Task{{ID: 1, Text: "Item to check", Done: false}},
			want:    []Task{{ID: 1, Text: "Item to check", Done: false}},
			wantErr: ErrInvalidID,
		},
		{
			name:    "failed, task already completed",
			input:   "1",
			list:    []Task{{ID: 1, Text: "Item to check", Done: true}},
			want:    []Task{{ID: 1, Text: "Item to check", Done: true}},
			wantErr: ErrAlreadyDone,
		},
		{
			name:    "failed, no tasks in list",
			input:   "1",
			list:    []Task{},
			want:    []Task{},
			wantErr: ErrNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := &mockStore{tasks: test.list}
			f := NewTaskService(store)

			err := f.Check(test.input)
			if err != nil && !errors.Is(err, test.wantErr) {
				t.Fatalf("got %v, expected %v", err, test.wantErr)
			}
			if !reflect.DeepEqual(store.tasks, test.want) {
				t.Fatalf("got %+v, expected %+v", store.tasks, test.want)
			}
		})
	}
}
