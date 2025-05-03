package tgapi

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkrtychanr/rag_bot/internal/model"
)

type tgAPI struct {
	gateway *tgbotapi.BotAPI
}

func NewTgAPI(api *tgbotapi.BotAPI) *tgAPI {
	return &tgAPI{
		gateway: api,
	}
}

func (api *tgAPI) SendScreen(chatID int64, screen model.Screen) (int64, error) {
	msg := tgbotapi.NewMessage(chatID, screen.Text)

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(screen.Buttons))

	for _, buttons := range screen.Buttons {
		row := make([]tgbotapi.InlineKeyboardButton, 0, len(buttons))

		for _, button := range buttons {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(button.Text, string(button.Payload)))
		}

		rows = append(rows, row)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	msg.ReplyMarkup = &keyboard

	m, err := api.gateway.Send(msg)
	if err != nil {
		return 0, fmt.Errorf("failed to send message: %w", err)
	}

	return int64(m.MessageID), nil
}

func (api *tgAPI) DeleteMessage(chatID int64, messageID int64) error {
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, int(messageID))

	if _, err := api.gateway.Request(deleteMsg); err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}

	return nil
}

func (api *tgAPI) SendText(chatID int64, text string) error {
	if _, err := api.gateway.Send(tgbotapi.NewMessage(chatID, text)); err != nil {
		return fmt.Errorf("failed to send: %w", err)
	}

	return nil
}

func (api *tgAPI) GetUserShortname(tgID int64) (string, error) {
	result, err := api.gateway.GetChat(tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: tgID,
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to get chat: %w", err)
	}

	return result.UserName, nil
}
