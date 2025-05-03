package screenchanger

import (
	"context"
	"errors"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/model"
	"github.com/mkrtychanr/rag_bot/internal/repository"
)

type useCase struct {
	gateway gateway
	repo    repo
}

func NewUseCase(g gateway, r repo) *useCase {
	return &useCase{
		gateway: g,
		repo:    r,
	}
}

func (u *useCase) ChangeScreen(ctx context.Context, chatID int64, tgID int64, newScreen model.Screen) error {
	result, err := u.repo.GetUserIDByTelegramID(ctx, tgID)
	if errors.Is(err, repository.ErrNoRows) {
		shortName, err := u.gateway.GetUserShortname(tgID)
		if err != nil {
			return fmt.Errorf("failed to get user shortname: %w", err)
		}

		userID, err := u.repo.RegistrateUser(ctx, shortName, tgID)
		if err != nil {
			return fmt.Errorf("failed to registrate user: %w", err)
		}

		result = userID
	}
	if err != nil {
		return fmt.Errorf("failed to get user id by telegram id: %w", err)
	}

	newMessageID, err := u.gateway.SendScreen(chatID, newScreen)
	if err != nil {
		return fmt.Errorf("failed to send screen: %w", err)
	}

	if err := u.repo.ChangeMessage(ctx, chatID, result, newMessageID); err != nil {
		return fmt.Errorf("failed to change message: %w", err)
	}

	return nil
}
