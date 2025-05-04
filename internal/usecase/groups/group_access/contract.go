package groupaccess

import (
	"context"

	"github.com/mkrtychanr/rag_bot/internal/model"
)

type repository interface {
	GetUserGroups(ctx context.Context, tgID int64) ([]int64, error)
	GetUserRWGroupIDs(ctx context.Context, tgID int64) ([]int64, error)
	GetUserGroupsOwnership(ctx context.Context, tgID int64) ([]model.Group, error)
	GetPaperGroupIDs(ctx context.Context, paperID int64) ([]int64, error)
	FetchGroupsInfo(ctx context.Context, groupIDs []int64) ([]model.Group, error)
}
