package storage

import (
	"io"
	"os"
	"path/filepath"
)

type Storage interface {
	Save(name string, file *os.File) error
	GetByName(name string) (*os.File, error)
}

type PhotoStorage struct {
	path string
}

func NewPhotoStorage(path string) *PhotoStorage {
	return &PhotoStorage{
		path: path,
	}
}

func (ps *PhotoStorage) Save(name string, file *os.File) error {
	fullPath := filepath.Join(ps.path, name)

	outFile, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PhotoStorage) GetByName(name string) (*os.File, error) {
	fullPath := filepath.Join(ps.path, name)

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}

	return file, nil
}
