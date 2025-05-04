package deletedocumentfromgroupscreen

import (
	"context"

	"github.com/mkrtychanr/rag_bot/internal/model"
)

type deleter interface {
	DeleteDocumentFromGroup(ctx context.Context, paperID int64, groupID int64) error
}

type getter interface {
	GetGroupsToDeletePaper(ctx context.Context, tgID int64, paperID int64) ([]model.Group, error)
}
