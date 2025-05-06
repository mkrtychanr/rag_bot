package makerequest

import (
	"context"

	"github.com/mkrtychanr/rag_bot/internal/model"
)

type repository interface {
	GetUserGroups(ctx context.Context, tgID int64) ([]int64, error)
	GetGroupPapers(ctx context.Context, groupID int64) ([]int64, error)
	GetUserGroupsOwnership(ctx context.Context, tgID int64) ([]model.Group, error)
}

type rag interface {
	GetLLMResponse(ctx context.Context, request string, paperIDs []int64) (string, error)
}
