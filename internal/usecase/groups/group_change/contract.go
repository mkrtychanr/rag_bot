package groupchange

import (
	"context"

	"github.com/mkrtychanr/rag_bot/internal/model"
)

type repository interface {
	ChangeGroupName(ctx context.Context, groupID int64, name string) error
	AddUserIntoGroup(ctx context.Context, groupID int64, tgID int64) error
	DeleteUserFromGroup(ctx context.Context, groupID int64, tgID int64) error
	SetReadOnlyRightsForUserInGroup(ctx context.Context, groupID int64, tgID int64) error
	SetReadWriteRightsForUserInGroup(ctx context.Context, groupID int64, tgID int64) error
	GetUserIDsInGroupWithReadOnlyAccess(ctx context.Context, groupID int64) ([]int64, error)
	GetUserIDsInGroupWithReadWriteAccess(ctx context.Context, groupID int64) ([]int64, error)
	FetchUsersInfo(ctx context.Context, userIDs []int64) ([]model.User, error)
	GetUserByShortname(ctx context.Context, shortname string) (model.User, error)
	GetGroupUsers(ctx context.Context, groupID int64) ([]model.UserGroup, error)
}
