package makerequest

import "context"

type repository interface {
	GetUserGroups(ctx context.Context, tgID int64) ([]int64, error)
	GetGroupPapers(ctx context.Context, groupID int64) ([]int64, error)
}

type rag interface {
	GetLLMResponse(ctx context.Context, request string, papers []int64) (string, error)
}
