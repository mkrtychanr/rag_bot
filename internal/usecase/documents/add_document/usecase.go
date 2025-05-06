package adddocument

import (
	"context"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/logger"
)

type useCase struct {
	getter     getter
	uploader   uploader
	repository repository
}

func NewUseCase(g getter, u uploader, r repository) *useCase {
	return &useCase{
		getter:     g,
		uploader:   u,
		repository: r,
	}
}

func (u *useCase) AddDocument(ctx context.Context, name string, tgID int64, documentTelegramID string) error {
	file, err := u.getter.GetFile(ctx, documentTelegramID)
	if err != nil {
		return fmt.Errorf("failed to get file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			logger.GetLogger().Err(err).Msg("failed to close file")
		}
	}()

	documentID, err := u.repository.AddPaper(ctx, name, tgID)
	if err != nil {
		return fmt.Errorf("failed to add paper into repository: %w", err)
	}

	if err := u.uploader.UploadDocument(ctx, file, documentID); err != nil {
		return fmt.Errorf("failed to upload document into rag: %w", err)
	}

	return nil
}
