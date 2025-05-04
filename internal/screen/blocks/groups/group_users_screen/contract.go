package groupusersscreen

import (
	"context"

	"github.com/mkrtychanr/rag_bot/internal/model"
)

type getter interface {
	GetGroupUsers(ctx context.Context, groupID int64) ([]model.UserGroup, error)
}
