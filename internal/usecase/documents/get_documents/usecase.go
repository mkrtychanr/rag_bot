package getdocuments

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

func (u *useCase) GetUserDocuments(ctx context.Context, tgID int64) ([]model.Paper, error) {
	papers, err := u.repository.GetUserPapers(ctx, tgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user papers: %w", err)
	}

	return papers, nil
}
