package pickinfogroup

import (
	"context"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/model"
	"github.com/mkrtychanr/rag_bot/internal/utils"
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

	m := make(map[int64]model.Group, len(accessGroups)+len(ownedGroups))

	for _, group := range accessGroups {
		m[group.ID] = group
	}

	for _, group := range ownedGroups {
		m[group.ID] = group
	}

	return utils.MapValuesToSlice(m), nil
}
