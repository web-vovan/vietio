package ads

import (
	"context"
	"mime/multipart"
)

type FileStorage interface {
    Save(ctx context.Context, file multipart.File, header *multipart.FileHeader) (FileInfo, error)
    DeleteByPath(ctx context.Context, path string) error
}

type FileInfo struct {
    FileName string
    Size int64
    Mime string
}