package groupownership

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

func (u *useCase) GetUserGroupsOwnership(ctx context.Context, tgID int64) ([]model.Group, error) {
	groups, err := u.repository.GetUserGroupsOwnership(ctx, tgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups ownership: %w", err)
	}

	return groups, nil
}
