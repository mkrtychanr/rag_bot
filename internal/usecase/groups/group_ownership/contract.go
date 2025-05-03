package groupownership

import (
	"context"

	"github.com/mkrtychanr/rag_bot/internal/model"
)

type repository interface {
	GetUserGroupsOwnership(ctx context.Context, tgID int64) ([]model.Group, error)
}
