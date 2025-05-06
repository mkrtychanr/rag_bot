package adddocument

import (
	"context"
	"io"
)

type getter interface {
	GetFile(ctx context.Context, fileID string) (io.ReadCloser, error)
}

type repository interface {
	AddPaper(ctx context.Context, name string, tgID int64) (int64, error)
}

type uploader interface {
	UploadDocument(ctx context.Context, file io.ReadCloser, id int64) error
}
