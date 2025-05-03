package changedocument

import "context"

type repository interface {
	AddPaperIntoGroup(ctx context.Context, paperID int64, groupID int64) error
	DeletePaperFromGroup(ctx context.Context, paperID int64, groupID int64) error
}
