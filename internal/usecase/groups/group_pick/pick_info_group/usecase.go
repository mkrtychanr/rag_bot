package pickinfogroup

import (
	"context"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/model"
)

type useCase struct {
	repository repository
}

func NewUseCase(r repository) *useCase {
	return &useCase{
		repository: r,
	}
}

func (u *useCase) GetGroups(ctx context.Context, tgID int64) ([]model.Group, error) {
	hasAccess, err := u.repository.GetUserGroups(ctx, tgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups: %w", err)
	}

	accessGroups, err := u.repository.FetchGroupsInfo(ctx, hasAccess)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch groups info: %w", err)
	}

	ownedGroups, err := u.repository.GetUserGroupsOwnership(ctx, tgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups ownership: %w", err)
	}

	result := make([]model.Group, 0, len(accessGroups)+len(ownedGroups))

	result = append(result, accessGroups...)
	result = append(result, ownedGroups...)

	return result, nil
}
