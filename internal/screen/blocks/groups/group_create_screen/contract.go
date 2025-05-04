package groupcreatescreen

import "context"

type creator interface {
	CreateGroup(ctx context.Context, tgID int64, name string) error
}
