package groupadduserscreen

import "context"

type adder interface {
	AddUserIntoGroup(ctx context.Context, groupID int64, name string) error
}
