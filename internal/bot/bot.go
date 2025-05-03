package bot

import (
	"context"
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkrtychanr/rag_bot/internal/config"
	"github.com/mkrtychanr/rag_bot/internal/logger"
)

var ErrUpdateChanClosed = errors.New("udpates chan is closed")

type impl struct {
	reciver messageReciver
	offset  int
	timeout int
}

func NewBot(reciver messageReciver, cfg config.Bot) (*impl, error) {
	return &impl{
		reciver: reciver,
		offset:  cfg.Offset,
		timeout: cfg.Timeout,
	}, nil
}

// var (
// 	first bool
// 	msgID int
// )

func (i *impl) Recive(ctx context.Context, handleChan chan<- tgbotapi.Update) error {
	cfg := tgbotapi.NewUpdate(i.offset)
	cfg.Timeout = i.timeout

	updates := i.reciver.GetUpdatesChan(cfg)

	logger.GetLogger().Info().Msg("bot run started")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case update, ok := <-updates:
			if !ok {
				return ErrUpdateChanClosed
			}

			// logger.GetLogger().Info().Msgf("%v", update)
			// if update.CallbackQuery != nil {
			// 	logger.GetLogger().Info().Msgf("%v", update.CallbackQuery)
			// }

			// if update.Message != nil {
			// 	if update.Message.Text == "" {
			// 		continue
			// 	}

			// 	logger.GetLogger().Info().Msg("new message")
			// 	id := update.Message.Chat.ID

			// 	if !first {
			// 		button := tgbotapi.NewInlineKeyboardButtonData("Кнопка", `{"data": "aboba"}`)

			// 		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			// 			tgbotapi.NewInlineKeyboardRow(button),
			// 		)

			// 		// Создаем сообщение с кнопками
			// 		msg := tgbotapi.NewMessage(id, update.Message.Text)
			// 		msg.ReplyMarkup = keyboard

			// 		m, err := i.reciver.Send(msg)
			// 		if err != nil {
			// 			logger.GetLogger().Err(err).Msg("failed to send message")
			// 		}
			// 		first = true
			// 		msgID = m.MessageID

			// 		continue
			// 	}

			// 	deleteMsg := tgbotapi.NewDeleteMessage(id, msgID)
			// 	if _, err := i.reciver.Request(deleteMsg); err != nil {
			// 		logger.GetLogger().Err(err).Msg("failed to delete message")
			// 	}

			// 	newMsg := tgbotapi.NewMessage(id, update.Message.Text)
			// 	button := tgbotapi.NewInlineKeyboardButtonData("Кнопка", `{"data": "biba"}`)

			// 	keyboard := tgbotapi.NewInlineKeyboardMarkup(
			// 		tgbotapi.NewInlineKeyboardRow(button),
			// 	)
			// 	newMsg.ReplyMarkup = &keyboard
			// 	m, err := i.reciver.Send(newMsg)
			// 	if err != nil {
			// 		logger.GetLogger().Err(err).Msg("failed to send message")
			// 	}

			// 	msgID = m.MessageID

			handleChan <- update

		}
	}
}

// func newBot(cfg config.Bot) (*impl, error) {
// 	b, err := tgbotapi.NewBotAPI(cfg.Token)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create bot instance: %w", err)
// 	}

// 	updateCfg := atgbotapiNewUpdate(0)

// 	updateCfg.Timeout = 5

// 	updates := b.GetUpdatesChan(updateCfg)

// 	for update := range updates {
// 		if update.Message == nil {
// 			continue
// 		}

// 		if update.Message.Text == "set commands" {
// 			r, err := b.MakeRequest("setMyCommands", map[string]string{"commands": `[{"command": "aboba", "description": "biba"}]`})
// 			if err != nil {
// 				return nil, fmt.Errorf("failed to set my commands: %w", err)
// 			}

// 			fmt.Println(r)
// 		}
// 	}

// 	return nil, nil
// }
