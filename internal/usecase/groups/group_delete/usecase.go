package groupdelete

import (
	"context"
	"fmt"
)

type useCase struct {
	repository repository
}

func NewUseCase(repository repository) *useCase {
	return &useCase{
		repository: repository,
	}
}

func (u *useCase) DeleteGroup(ctx context.Context, groupID int64) error {
	if err := u.repository.DeleteGroups(ctx, []int64{groupID}); err != nil {
		return fmt.Errorf("failed to delete groups: %w", err)
	}

	return nil
}
