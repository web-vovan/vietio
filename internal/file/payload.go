package file

import "github.com/google/uuid"

type File struct {
    Id int64
    AdUuid uuid.UUID
    Path string
    PreviewPath string
    Order int
    Size int64
    PreviewSize int64
    Mime string
    PreviewMime string
}