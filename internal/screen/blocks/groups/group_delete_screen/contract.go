package groupdeletescreen

import (
	"context"

	"github.com/mkrtychanr/rag_bot/internal/model"
)

type getter interface {
	GetGroups(ctx context.Context, tgID int64) ([]model.Group, error)
}

type deleter interface {
	DeleteGroup(ctx context.Context, groupID int64) error
}
