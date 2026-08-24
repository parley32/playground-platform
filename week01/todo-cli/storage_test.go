package main

import (
	"os"
	"reflect"
	"testing"
)

func TestStoreAndLoadTodos(t *testing.T) {
	t.Chdir(t.TempDir())

	want := []Todo{
		{Name: "buy milk", ID: 1, Completed: false},
		{Name: "walk the dog", ID: 2, Completed: true},
	}

	if err := storeTodos(&want); err != nil {
		t.Fatalf("storeTodos() error = %v", err)
	}

	file, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("open stored todos: %v", err)
	}
	t.Cleanup(func() { file.Close() })

	got, err := loadTodos(file)
	if err != nil {
		t.Fatalf("loadTodos() error = %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("loadTodos() = %#v, want %#v", got, want)
	}
}

func TestLoadTodosEmptyFile(t *testing.T) {
	t.Chdir(t.TempDir())

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		t.Fatalf("create empty todos file: %v", err)
	}
	t.Cleanup(func() { file.Close() })

	got, err := loadTodos(file)
	if err != nil {
		t.Fatalf("loadTodos() error = %v", err)
	}

	if got == nil {
		t.Fatal("loadTodos() returned nil slice for empty file")
	}
	if len(got) != 0 {
		t.Errorf("loadTodos() returned %d todos, want 0", len(got))
	}
}
