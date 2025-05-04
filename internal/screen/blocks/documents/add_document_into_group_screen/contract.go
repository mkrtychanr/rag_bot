package adddocumentintogroupscreen

import (
	"context"

	"github.com/mkrtychanr/rag_bot/internal/model"
)

type adder interface {
	AddDocumentIntoGroup(ctx context.Context, paperID int64, groupID int64) error
}

type getter interface {
	GetGroupsToAddPaper(ctx context.Context, tgID int64, paperID int64) ([]model.Group, error)
}
