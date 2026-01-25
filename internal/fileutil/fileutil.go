package fileutil

import (
	"errors"
	"os"
	"path/filepath"
)

type FileInfo struct {
	Path string
	Size int64
}

func GetFileInfo(path string) (*FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, errors.New("path is a directory")
	}

	return &FileInfo{
		Path: path,
		Size: info.Size(),
	}, nil
}

func ListFiles(dir string, recursive bool) ([]FileInfo, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []FileInfo
	for _, e := range entries {
		fullPath := filepath.Join(dir, e.Name())
		if e.IsDir() {
			if recursive {
				subFiles, err := ListFiles(fullPath, true)
				if err != nil {
					continue
				}
				files = append(files, subFiles...)
			}
			continue
		}

		info, err := e.Info()
		if err != nil {
			continue
		}

		files = append(files, FileInfo{
			Path: fullPath,
			Size: info.Size(),
		})
	}

	return files, nil
}
