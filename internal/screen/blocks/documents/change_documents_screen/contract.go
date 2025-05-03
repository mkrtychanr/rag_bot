package changedocumentsscreen

import (
	"context"

	"github.com/mkrtychanr/rag_bot/internal/model"
)

type getter interface {
	GetUserDocuments(ctx context.Context, tgID int64) ([]model.Paper, error)
}
