package deletedocument

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

func (u *useCase) DeleteDocument(ctx context.Context, paperID int64) error {
	if err := u.repository.DeletePapers(ctx, []int64{paperID}); err != nil {
		return fmt.Errorf("failed to delete paper: %w", err)
	}

	return nil
}
