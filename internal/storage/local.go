package storage

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"vietio/internal/ads"

	"github.com/google/uuid"
)

type LocalStorage struct {
	BasePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{
		BasePath: basePath,
	}
}

func (s *LocalStorage) Save(
	ctx context.Context,
	file multipart.File,
	header *multipart.FileHeader,
) (ads.FileInfo, error) {
	ext := filepath.Ext(header.Filename)
	fileName := uuid.NewString() + ext

	path := filepath.Join(s.BasePath, fileName)

	dst, err := os.Create(path)
	if err != nil {
		return ads.FileInfo{}, err
	}
	defer dst.Close()

	size, err := io.Copy(dst, file)
	if err != nil {
		return ads.FileInfo{}, err
	}

	return ads.FileInfo{
		FileName: fileName,
		Size: size,
		Mime: header.Header.Get("Content-Type"),
	}, nil
}
