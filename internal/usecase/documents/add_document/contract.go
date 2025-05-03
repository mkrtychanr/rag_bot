package adddocument

import "context"

type repository interface {
	AddPaper(ctx context.Context, name string, tgID int64) (int64, error)
}

type rag interface {
	UploadDocument(ctx context.Context, documentID int64) error
}
