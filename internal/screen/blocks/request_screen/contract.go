package requestscreen

import "context"

type requestMaker interface {
	MakeRequest(ctx context.Context, request string, tgID int64) (string, error)
}

type sender interface {
	SendText(ctx context.Context, userID int64, message string) error
}
