package groupaccess

import (
	"context"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/model"
)

type useCase struct {
	repository repository
}

func NewUseCase(repository repository) *useCase {
	return &useCase{
		repository: repository,
	}
}

func (u *useCase) GetUserAccessGroups(ctx context.Context, tgID int64) ([]model.Group, error) {
	groupIDs, err := u.repository.GetUserGroups(ctx, tgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups: %w", err)
	}

	groups, err := u.repository.FetchGroupsInfo(ctx, groupIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch groups info: %w", err)
	}

	return groups, nil
}
