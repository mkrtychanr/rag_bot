package changedocument

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

func (u *useCase) AddDocumentIntoGroup(ctx context.Context, paperID int64, groupID int64) error {
	if err := u.repository.AddPaperIntoGroup(ctx, paperID, groupID); err != nil {
		return fmt.Errorf("failed to add paper into group: %w", err)
	}

	return nil
}

func (u *useCase) DeleteDocumentFromGroup(ctx context.Context, paperID int64, groupID int64) error {
	if err := u.repository.DeletePaperFromGroup(ctx, paperID, groupID); err != nil {
		return fmt.Errorf("failed to delete paper from group: %w", err)
	}

	return nil
}
