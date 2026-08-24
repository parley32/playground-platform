package main

import (
	"os"
	"reflect"
	"slices"
	"testing"
)

func TestGetHighestCurrentID(t *testing.T) {
	tests := []struct {
		name   string
		idList []Todo
		wantID int
	}{
		{
			name:   "returns highest ID",
			idList: []Todo{{ID: 0}, {ID: 1}, {ID: 2}},
			wantID: 2,
		},
		{
			name:   "empty list",
			idList: []Todo{},
			wantID: 0,
		},
		{
			name:   "IDs are out of order",
			idList: []Todo{{ID: 3}, {ID: 1}},
			wantID: 3,
		},
		{
			name:   "zero ID",
			idList: []Todo{{ID: 0}},
			wantID: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := getHighestCurrentID(test.idList)
			if got != test.wantID {
				t.Errorf("getHighestCurrentID() = %d, want %d", got, test.wantID)
			}
		})
	}
}

func TestAddItemToList(t *testing.T) {
	t.Chdir(t.TempDir())

	want := []Todo{
		{Name: "existing task", ID: 1},
		{Name: "another task", ID: 4},
	}
	if err := storeTodos(&want); err != nil {
		t.Fatalf("storeTodos() error = %v", err)
	}

	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		t.Fatalf("open todos file: %v", err)
	}
	t.Cleanup(func() { file.Close() })

	gotMessage, err := addItemToList(file, "new task")
	if err != nil {
		t.Fatalf("addItemToList() error = %v", err)
	}
	if gotMessage == "" {
		t.Fatal("addItemToList() returned an empty success message")
	}

	got, err := loadTodos(file)
	if err != nil {
		t.Fatalf("loadTodos() error = %v", err)
	}

	want = append(want, Todo{Name: "new task", ID: 5})
	if !reflect.DeepEqual(got, want) {
		t.Errorf("loadTodos() after add = %#v, want %#v", got, want)
	}
}

func TestRemoveItemFromList(t *testing.T) {
	t.Chdir(t.TempDir())

	want := []Todo{
		{Name: "existing task", ID: 1},
		{Name: "another task", ID: 4},
	}
	if err := storeTodos(&want); err != nil {
		t.Fatalf("storeTodos() error = %v", err)
	}

	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		t.Fatalf("open todos file: %v", err)
	}
	t.Cleanup(func() { file.Close() })

	gotMessage, err := removeItemFromList(file, 1)
	if err != nil {
		t.Fatalf("removeItemFromList() error = %v", err)
	}
	if gotMessage == "" {
		t.Fatal("removeItemFromList() returned an empty success message")
	}

	got, err := loadTodos(file)
	if err != nil {
		t.Fatalf("loadTodos() error = %v", err)
	}

	want = []Todo{
		{Name: "another task", ID: 4},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("loadTodos() after add = %#v, want %#v", got, want)
	}
}

func TestCheckItem(t *testing.T) {
	tests := []struct {
		name    string
		todo    []Todo
		want    []Todo
		wantErr bool
	}{
		{
			name: "checking a non-complete task",
			todo: []Todo{{Name: "test1", ID: 1, Completed: false}},
			want: []Todo{{Name: "test1", ID: 1, Completed: true}},
		},
		{
			name:    "checking an already completed task",
			todo:    []Todo{{Name: "test2", ID: 1, Completed: true}},
			want:    []Todo{{Name: "test2", ID: 1, Completed: true}},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Chdir(t.TempDir())

			input := test.todo
			if err := storeTodos(&input); err != nil {
				t.Fatalf("storeTodos() error = %v", err)
			}

			file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
			if err != nil {
				t.Fatalf("open todos file: %v", err)
			}
			t.Cleanup(func() { file.Close() })

			gotMessage, err := checkItem(file, 1)
			if (err != nil) != test.wantErr {
				t.Fatalf("checkItem() error = %v", err)
			}
			if !test.wantErr && gotMessage == "" {
				t.Fatal("checkItem() returned an empty success message")
			}
			if test.wantErr && gotMessage != "" {
				t.Errorf("checkItem() message = %q, want empty message", gotMessage)
			}

			got, err := loadTodos(file)
			if err != nil {
				t.Fatalf("loadTodos() error = %v", err)
			}

			if !slices.Equal(got, test.want) {
				t.Errorf("loadTodos() after done = %#v, want %#v", got, test.want)
			}
		})
	}
}
