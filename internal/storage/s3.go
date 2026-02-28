package storage

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"vietio/internal/ads"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

type S3Storage struct {
	client     *s3.Client
	bucketName string
	publicURL  string
}

func NewS3Storage(
	ctx context.Context,
	key, secret, bucket, publicURL string,
) (*S3Storage, error) {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion("ru-central1"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(key, secret, ""),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("https://storage.yandexcloud.net")
		o.UsePathStyle = true //
	})

	return &S3Storage{
		client:     client,
		bucketName: bucket,
		publicURL:  publicURL,
	}, nil
}

func (s *S3Storage) Save(
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
	if img.Bounds().Dx() > 1000 {
		fullImg = imaging.Resize(img, 1000, 0, imaging.Lanczos)
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

	// отправка изображения в s3
    fileUUID := uuid.NewString()
	fullFileName := fileUUID + ".jpg"
	if err := s.uploadObject(ctx, fullFileName, fullImgData); err != nil {
		return nil, fmt.Errorf("failed to upload full image: %w", err)
	}

	// отправка превью в s3
    previewFileName := fileUUID + "_preview.jpg"
	if err := s.uploadObject(ctx, previewFileName, previewImgData); err != nil {
		return nil, fmt.Errorf("failed to upload preview image: %w", err)
	}

	return &ads.FileInfo{
		FileName:        fullFileName,
		PreviewFileName: previewFileName,
		Size:            int64(len(fullImgData)),
		PreviewSize:     int64(len(previewImgData)),
		Mime:            "image/jpeg",
		PreviewMime:     "image/jpeg",
	}, nil
}

// uploadObject — вспомогательный метод для загрузки байтов в S3
func (s *S3Storage) uploadObject(ctx context.Context, key string, data []byte) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data), // Превращаем []byte в io.Reader
		ContentType: aws.String("image/jpeg"),
		ACL:         types.ObjectCannedACLPublicRead, // Делаем файл публично доступным для чтения
	})
	return err
}

func (s *S3Storage) DeleteByPath(ctx context.Context, path string) error {	
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}
	return nil
}

func (s *S3Storage) GetPublicPath(path string) string {
	// Просто склеиваем базовый урл бакета и имя файла
	return fmt.Sprintf("%s/%s", s.publicURL, path)
}

func (s *S3Storage) GetType() string {
	return "s3"
}