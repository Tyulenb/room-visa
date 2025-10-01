package storage

import (
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

)

type Storage interface {
	Save(name string, file multipart.File) error
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

func (ps *PhotoStorage) Save(name string, file multipart.File) error {
    img, format, err := image.Decode(file) 
    if err != nil {
        return err
    }

	fullPath := filepath.Join(ps.path, name)

	outFile, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

    switch strings.ToLower(format) {
    case "jpeg", "jpg":
        if err := jpeg.Encode(outFile, img, nil); err != nil {
            return err
        }
    case "png":
        if err := png.Encode(outFile, img); err != nil {
            return err
        }
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
