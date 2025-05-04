package screenchanger

import (
	"context"

	"github.com/mkrtychanr/rag_bot/internal/model"
)

type repo interface {
	GetUserIDByTelegramID(ctx context.Context, tgID int64) (int64, error)
	RegistrateUser(ctx context.Context, name string, tgID int64) (int64, error)
	ChangeMessage(ctx context.Context, chatID int64, userID int64, messageID int64) error
	GetMessageID(ctx context.Context, chatID int64) (int64, error)
	DeleteMessage(ctx context.Context, chatID int64) error
}

type gateway interface {
	GetUserShortname(ctx context.Context, tgID int64) (string, error)
	SendScreen(ctx context.Context, chatID int64, screen model.Screen) (int64, error)
	DeleteMessage(ctx context.Context, chatID int64, messageID int64) error
}
