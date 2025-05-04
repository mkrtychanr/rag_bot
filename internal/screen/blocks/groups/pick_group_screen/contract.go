package pickgroupscreen

import (
	"context"

	"github.com/mkrtychanr/rag_bot/internal/model"
)

type getter interface {
	GetGroups(ctx context.Context, tgID int64) ([]model.Group, error)
}

type performer interface {
	Perform(ctx context.Context, payload map[string]any) error
}
