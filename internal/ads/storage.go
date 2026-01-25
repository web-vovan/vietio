package ads

import (
	"context"
	"mime/multipart"
)

type FileStorage interface {
    Save(ctx context.Context, file multipart.File, header *multipart.FileHeader) (*FileInfo, error)
    DeleteByPath(ctx context.Context, path string) error
    GetPublicPath(path string) string
}

type FileInfo struct {
    FileName string
    PreviewFileName string
    Size int64
    PreviewSize int64
    Mime string
    PreviewMime string
}