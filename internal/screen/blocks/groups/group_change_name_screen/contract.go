package groupchangenamescreen

import "context"

type changer interface {
	ChangeGroupName(ctx context.Context, groupID int64, name string) error
}
