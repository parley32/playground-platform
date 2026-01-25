package fileutil_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/parley-32/playground-platform/internal/fileutil"
)

func TestListFiles_NonRecursive(t *testing.T) {
	dir := t.TempDir()

	// create files
	os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(dir, "b.txt"), []byte("b"), 0644)

	// create subdirectory with file
	sub := filepath.Join(dir, "sub")
	os.Mkdir(sub, 0755)
	os.WriteFile(filepath.Join(sub, "c.txt"), []byte("c"), 0644)

	files, err := fileutil.ListFiles(dir, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}

}

func TestListFiles_Recursive(t *testing.T) {
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644)

	sub := filepath.Join(dir, "sub")
	os.Mkdir(sub, 0755)
	os.WriteFile(filepath.Join(sub, "b.txt"), []byte("b"), 0644)

	files, err := fileutil.ListFiles(dir, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}

}

func TestListFiles_InvalidDir(t *testing.T) {
	_, err := fileutil.ListFiles("/definitely/not/a/real/path", false)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetFileInfo_ValidFile(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.txt")
	content := []byte("hello world")

	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	fi, err := fileutil.GetFileInfo(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fi.Path != filePath {
		t.Fatalf("expected path %q, got %q", filePath, fi.Path)
	}

	if fi.Size != int64(len(content)) {
		t.Fatalf("expected size %d, got %d", len(content), fi.Size)
	}
}

func TestGetFileInfo_Directory(t *testing.T) {
	dir := t.TempDir()

	_, err := fileutil.GetFileInfo(dir)
	if err == nil {
		t.Fatal("expected error for directory path, got nil")
	}
}

func TestGetFileInfo_NonExistentFile(t *testing.T) {
	_, err := fileutil.GetFileInfo("/definitely/not/a/real/file.txt")
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
}
