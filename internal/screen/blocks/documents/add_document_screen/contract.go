package adddocumentscreen

import "context"

type documentAdder interface {
	AddDocument(ctx context.Context, name string, tgID int64, documentTelegramID string) error
}
