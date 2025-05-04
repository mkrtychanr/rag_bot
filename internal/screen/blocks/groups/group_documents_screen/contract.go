package groupdocumentsscreen

import (
	"context"

	"github.com/mkrtychanr/rag_bot/internal/model"
)

type getter interface {
	GetGroupDocuments(ctx context.Context, groupID int64) ([]model.Paper, error)
}
