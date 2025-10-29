package storage

import (
	"os"
	"path/filepath"
)

type FileSystemStorage struct {
	workingDir string
}

func NewFileSystemStorage(workingDir string) *FileSystemStorage {
	return &FileSystemStorage{
		workingDir: workingDir,
	}
}

func (fs *FileSystemStorage) Store(filename string, data []byte) (string, error) {
	fullPath := filepath.Join(fs.workingDir, filename)

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	err := os.WriteFile(fullPath, data, 0644)
	if err != nil {
		return "", err
	}

	return fullPath, nil
}

func (fs *FileSystemStorage) Load(filename string) ([]byte, error) {
	fullPath := filepath.Join(fs.workingDir)

	return os.ReadFile(fullPath)
}
