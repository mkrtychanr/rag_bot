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
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			shortName, err := u.gateway.GetUserShortname(ctx, tgID)
			if err != nil {
				return fmt.Errorf("failed to get user shortname: %w", err)
			}

			userID, err := u.repo.RegistrateUser(ctx, shortName, tgID)
			if err != nil {
				return fmt.Errorf("failed to registrate user: %w", err)
			}

			result = userID
		} else {
			return fmt.Errorf("failed to get user id by telegram id: %w", err)
		}
	}

	if err := u.deleteMessage(ctx, chatID); err != nil {
		return err
	}

	newMessageID, err := u.gateway.SendScreen(ctx, chatID, newScreen)
	if err != nil {
		return fmt.Errorf("failed to send screen: %w", err)
	}

	if err := u.repo.ChangeMessage(ctx, chatID, result, newMessageID); err != nil {
		return fmt.Errorf("failed to change message: %w", err)
	}

	return nil
}

func (u *useCase) deleteMessage(ctx context.Context, chatID int64) error {
	messageID, err := u.repo.GetMessageID(ctx, chatID)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return nil
		}

		return fmt.Errorf("failed to get message id: %w", err)
	}

	if err := u.gateway.DeleteMessage(ctx, chatID, messageID); err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

func (u *useCase) Clear(ctx context.Context, chatID int64) error {
	_, err := u.repo.GetMessageID(ctx, chatID)
	if err != nil {
		if errors.Is(err, repository.ErrNoRows) {
			return nil
		}

		return fmt.Errorf("failed to get message id: %w", err)
	}

	if err := u.repo.DeleteMessage(ctx, chatID); err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}
