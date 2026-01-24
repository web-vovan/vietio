package storage

import (
	"context"
	"mime/multipart"
	"os"
	"path/filepath"
	"vietio/internal/ads"

	"github.com/disintegration/imaging"
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
) (*ads.FileInfo, error) {
	img, err := decodeImage(file)
	if err != nil {
		return nil, err
	}

	fileUUID := uuid.NewString()

	fullImg := img
	if img.Bounds().Dx() > 1200 {
		fullImg = imaging.Resize(img, 1200, 0, imaging.Lanczos)
	}

	fullFileName := fileUUID + ".jpg"
	fullPath := filepath.Join(s.BasePath, fullFileName)
	
	fullSize, err := saveAsJPG(fullImg, fullPath, 85)
	if err != nil {
		return nil, err
	}

	previewImg := imaging.Resize(fullImg, 300, 0, imaging.Lanczos)
	
	previewFileName := fileUUID + "_preview.jpg"
	previewPath := filepath.Join(s.BasePath, previewFileName)

	previewSize, err := saveAsJPG(previewImg, previewPath, 70)
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(header.Filename)
	fileName := uuid.NewString() + ext

	return &ads.FileInfo{
		FileName: fileName,
		PreviewFileName: previewFileName,
		Size: fullSize,
		PreviewSize: previewSize,
		Mime: "image/jpg",
		PreviewMime: "image/jpg",
	}, nil
}

func (s *LocalStorage) DeleteByPath(ctx context.Context, path string) error {
	deletePath := filepath.Join(s.BasePath, path)

	err := os.Remove(deletePath)
	if err != nil {
		return err
	}

	return nil
}
