package screenchanger

import (
	"context"

	"github.com/mkrtychanr/rag_bot/internal/model"
)

type repo interface {
	GetUserIDByTelegramID(ctx context.Context, tgID int64) (int64, error)
	RegistrateUser(ctx context.Context, name string, tgID int64) (int64, error)
	ChangeMessage(ctx context.Context, chat_id int64, user_id int64, message_id int64) error
}

type gateway interface {
	GetUserShortname(tgID int64) (string, error)
	SendScreen(chatID int64, screen model.Screen) (int64, error)
}
