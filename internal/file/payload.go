package file

import "github.com/google/uuid"

type FileModel struct {
	Id          int64
	AdUuid      uuid.UUID
	Path        string
	PreviewPath string
	Size        int64
	PreviewSize int64
	Mime        string
	PreviewMime string
	Storage     string
}
