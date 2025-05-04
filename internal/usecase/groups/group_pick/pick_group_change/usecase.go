package pickgroupchange

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
	ownedGroups, err := u.repository.GetUserGroupsOwnership(ctx, tgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups ownership: %w", err)
	}

	return ownedGroups, nil
}
