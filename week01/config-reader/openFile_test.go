package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenFile(t *testing.T) {
	t.Run("opens an existing file", func(t *testing.T) {
		directory := t.TempDir()
		fileName := filepath.Join(directory, "example.txt")

		if err := os.WriteFile(fileName, []byte("42"), 0o600); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		file, err := openFile(fileName)
		if err != nil {
			t.Fatalf("openFile() returned an unexpected error: %v", err)
		}
		defer file.Close()

		if file.Name() != fileName {
			t.Errorf("openFile() returned file %q, want %q", file.Name(), fileName)
		}
	})

	t.Run("returns an error for a missing file", func(t *testing.T) {
		file, err := openFile(filepath.Join(t.TempDir(), "missing.txt"))
		if err == nil {
			if file != nil {
				file.Close()
			}
			t.Fatal("openFile() returned nil error for a missing file")
		}

		if file != nil {
			file.Close()
			t.Fatal("openFile() returned a file for a missing path")
		}
	})
}
