package deletedocumentscreen

import (
	"context"

	"github.com/mkrtychanr/rag_bot/internal/model"
)

type deleter interface {
	DeleteDocument(ctx context.Context, paperID int64) error
}

type getter interface {
	GetUserDocuments(ctx context.Context, tgID int64) ([]model.Paper, error)
}
