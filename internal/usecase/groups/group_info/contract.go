package groupinfo

import (
	"context"

	"github.com/mkrtychanr/rag_bot/internal/model"
)

type repository interface {
	GetGroupPapers(ctx context.Context, groupID int64) ([]int64, error)
	GetGroupUsers(ctx context.Context, groupID int64) ([]model.UserGroup, error)
	FetchPapersInfo(ctx context.Context, paperIDs []int64) ([]model.Paper, error)
	FetchGroupsInfo(ctx context.Context, groupIDs []int64) ([]model.Group, error)
	FetchUsersInfo(ctx context.Context, userIDs []int64) ([]model.User, error)
}
