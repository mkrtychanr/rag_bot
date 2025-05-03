package groupcreate

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

func (u *useCase) CreateGroup(ctx context.Context, tgID int64, name string) error {
	if err := u.repository.CreateGroup(ctx, tgID, name); err != nil {
		return fmt.Errorf("failed to create group: %w", err)
	}

	return nil
}
