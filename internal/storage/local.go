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
	PublicUrl string
	BasePath string
}

func NewLocalStorage(publicUrl, basePath string) *LocalStorage {
	return &LocalStorage{
		PublicUrl: publicUrl,
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

	// изображение
	fullImg := img
	if img.Bounds().Dx() > 1200 {
		fullImg = imaging.Resize(img, 1200, 0, imaging.Lanczos)
	}

	fullImgData, err := encodeToJPG(fullImg, 85)
	if err != nil {
		return nil, err
	}

	// превью
	previewImg := imaging.Resize(fullImg, 300, 0, imaging.Lanczos)
	
	previewImgData, err := encodeToJPG(previewImg, 70)
	if err != nil {
		return nil, err
	}

	// сохраняем изображение
	fileUUID := uuid.NewString()
	fullFileName := fileUUID + ".jpg"
	fullPath := filepath.Join(s.BasePath, fullFileName)
	if err := os.WriteFile(fullPath, fullImgData, 0644); err != nil {
		return nil, err
	}

	// сохраняем превью
	previewFileName := fileUUID + "_preview.jpg"
	previewPath := filepath.Join(s.BasePath, previewFileName)
	if err := os.WriteFile(previewPath, previewImgData, 0644); err != nil {
		return nil, err
	}

	return &ads.FileInfo{
		FileName: fullFileName,
		PreviewFileName: previewFileName,
		Size: int64(len(fullImgData)),
		PreviewSize: int64(len(previewImgData)),
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

func (s *LocalStorage) GetPublicPath(path string) string {
	return s.PublicUrl + "/uploads/" + path
}

func (s *LocalStorage) GetType() string {
	return "local"
}