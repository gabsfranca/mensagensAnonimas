package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type StorageService interface {
	Save(file io.Reader, filename string) (string, error)
}

type LocalStorage struct {
	BasePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	os.MkdirAll(basePath, 0755)
	return &LocalStorage{BasePath: basePath}
}

func (s *LocalStorage) Save(file io.Reader, filename string) (string, error) {

	ext := filepath.Ext(filename)
	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	path := filepath.Join(s.BasePath, newFileName)

	dst, err := os.Create(path)
	if err != nil {
		return "", nil
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return "/media/" + newFileName, nil

}
