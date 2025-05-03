package deletedocument

import "context"

type repository interface {
	DeletePapers(ctx context.Context, papersIDs []int64) error
}
