package fileutil

import (
	"errors"
	"os"
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

func IsRegularFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.Mode().IsRegular(), nil
}