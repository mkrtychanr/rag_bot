package adddocument

import (
	"context"
	"fmt"
)

type useCase struct {
	rag        rag
	repository repository
}

func NewUseCase(rag rag, repository repository) *useCase {
	return &useCase{
		rag:        rag,
		repository: repository,
	}
}

func (u *useCase) AddDocument(ctx context.Context, name string, tgID int64, documentTelegramID string) error {
	documentID, err := u.repository.AddPaper(ctx, name, tgID)
	if err != nil {
		return fmt.Errorf("failed to add paper into repository: %w", err)
	}

	if err := u.rag.UploadDocument(ctx, documentID); err != nil {
		return fmt.Errorf("failed to upload document into rag: %w", err)
	}

	return nil
}
