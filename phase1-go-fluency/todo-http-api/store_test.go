package main

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestGetStore(t *testing.T) {
	tempDir := t.TempDir()

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		os.Chdir(oldDir)
	})

	contents := `[{"id":1,"text":"Buy milk","done":false}]`
	if err := os.WriteFile(fileName, []byte(contents), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := getStore()
	if err != nil {
		t.Fatal(err)
	}

	want := []Task{
		{ID: 1, Text: "Buy milk", Done: false},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestGetStoreMissingFile(t *testing.T) {
	tempDir := t.TempDir()

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		os.Chdir(oldDir)
	})

	got, err := getStore()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, []Task{}) {
		t.Fatalf("got %+v, want empty tasks", got)
	}
}

func TestSaveStore(t *testing.T) {
	tempDir := t.TempDir()

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		os.Chdir(oldDir)
	})

	tasks := []Task{{ID: 1, Text: "Buy Eggs", Done: false}, {ID: 2, Text: "Buys Milk", Done: false}}
	storedTasks, err := json.Marshal(tasks)
	if err != nil {
		t.Fatal(err)
	}
	err = saveStore(tasks)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, storedTasks) {
		t.Fatalf("got %q, expected %q", got, storedTasks)
	}

}
