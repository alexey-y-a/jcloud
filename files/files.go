package files

import (
	"os"
	"path/filepath"
	"strings"
)

type FileReader interface {
	ReadFile(filename string) ([]byte, error)
	IsJson(filename string) bool
}

type FileWriter interface {
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

type FileOperations interface {
	FileReader
	FileWriter
}

type FileSystem struct{}

func (fs FileSystem) ReadFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (fs FileSystem) IsJson(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".json"
}

func (fs FileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}
