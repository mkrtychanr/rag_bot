package groupdelete

import "context"

type repository interface {
	DeleteGroups(ctx context.Context, groupIDs []int64) error
}
