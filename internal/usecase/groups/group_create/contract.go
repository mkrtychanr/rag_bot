package groupcreate

import "context"

type repository interface {
	CreateGroup(ctx context.Context, tgID int64, name string) error
}
