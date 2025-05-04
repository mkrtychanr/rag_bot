package groupaccessscreen

import (
	"context"

	"github.com/mkrtychanr/rag_bot/internal/model"
)

type getter interface {
	GetUserAccessGroups(ctx context.Context, tgID int64) ([]model.Group, error)
}
