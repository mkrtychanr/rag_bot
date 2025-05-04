package smoothoperator

import (
	"context"
	"errors"
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mkrtychanr/rag_bot/internal/logger"
	"github.com/mkrtychanr/rag_bot/internal/model"
	"github.com/mkrtychanr/rag_bot/internal/screen"
)

var ErrUnknownChat = errors.New("unknown chat id")

const workersCount int64 = 20

type operator struct {
	getDefaultScreen func() screen.Screen
	chats            map[int64]screen.Screen
	chatsMu          sync.RWMutex
	screenChanger    screenChanger
	actionController actionController
}

func NewOperator(sc screenChanger, ac actionController, f func() screen.Screen) *operator {
	return &operator{
		screenChanger:    sc,
		getDefaultScreen: f,
		chats:            make(map[int64]screen.Screen),
		actionController: ac,
	}
}

func (o *operator) Operate(ctx context.Context, ch <-chan tgbotapi.Update) error {
	var wg sync.WaitGroup

	wg.Add(int(workersCount))

	for range workersCount {
		go func() {
			defer wg.Done()
			for {
				select {
				case v, ok := <-ch:
					if !ok {
						return
					}

					if err := o.handleUpdate(ctx, v); err != nil {
						logger.GetLogger().Err(err).Msg("failed to handle update")
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	wg.Wait()

	return nil
}

func (o *operator) handleUpdate(ctx context.Context, update tgbotapi.Update) error {
	if update.Message != nil {
		if update.Message.IsCommand() {
			if err := o.handleCommand(ctx, update); err != nil {
				return err
			}

			return nil
		}
	}

	if err := o.handleAction(ctx, getOperatorData(update)); err != nil {
		return fmt.Errorf("failed to handle action: %w", err)
	}

	return nil
}

func (o *operator) handleCommand(ctx context.Context, update tgbotapi.Update) error {
	cmd := update.Message.Command()

	switch cmd {
	case "start":
		if err := o.handleStartCommand(ctx, update.Message.Chat.ID, update.Message.From.ID); err != nil {
			return fmt.Errorf("failed to handle start command: %w", err)
		}
	}

	return nil
}

func (o *operator) handleStartCommand(ctx context.Context, chatID int64, tgID int64) error {
	o.chatsMu.Lock()
	defer o.chatsMu.Unlock()

	newScreen := o.getDefaultScreen()

	toChange, err := newScreen.Render()
	if err != nil {
		return fmt.Errorf("failed render new screen: %w", err)
	}

	if err := o.screenChanger.Clear(ctx, chatID); err != nil {
		return fmt.Errorf("failed to clear: %w", err)
	}

	if err := o.screenChanger.ChangeScreen(ctx, chatID, tgID, toChange); err != nil {
		return fmt.Errorf("failed to change screen: %w", err)
	}

	o.chats[chatID] = newScreen

	return nil
}

func getOperatorData(update tgbotapi.Update) model.OperatorData {
	res := model.OperatorData{}

	if update.CallbackQuery != nil {
		res = model.OperatorData{
			ChatID:       update.CallbackQuery.Message.Chat.ID,
			UserID:       update.CallbackQuery.From.ID,
			MessageID:    int64(update.CallbackQuery.Message.MessageID),
			CallbackData: []byte(update.CallbackQuery.Data),
		}
	} else if update.Message != nil {
		res = model.OperatorData{
			ChatID:    update.Message.Chat.ID,
			UserID:    update.Message.From.ID,
			MessageID: int64(update.Message.MessageID),
			Text:      &update.Message.Text,
		}

		if update.Message.Document != nil {
			res.DocumentID = &update.Message.Document.FileID
		}
	}

	return res
}

func (o *operator) handleAction(ctx context.Context, message model.OperatorData) error {
	curr, err := func() (screen.Screen, error) {
		o.chatsMu.RLock()
		defer o.chatsMu.RUnlock()

		curr, ok := o.chats[message.ChatID]
		if !ok {
			return nil, ErrUnknownChat
		}

		return curr, nil
	}()
	if err != nil {
		return err
	}

	f, err := o.actionController.GetScreenController(curr.GetScreenType())
	if err != nil {
		return fmt.Errorf("failed to get screen controller: %w", err)
	}

	curr, err = f(ctx, curr, message)
	if err != nil {
		return fmt.Errorf("failed to execute screen controller: %w", err)
	}

	toChange, err := curr.Render()
	if err != nil {
		return fmt.Errorf("failed to render new screen: %w", err)
	}

	if err := func() error {
		o.chatsMu.Lock()
		defer o.chatsMu.Unlock()

		if err := o.screenChanger.ChangeScreen(ctx, message.ChatID, message.UserID, toChange); err != nil {
			return fmt.Errorf("failed to change screen: %w", err)
		}

		o.chats[message.ChatID] = curr

		return nil
	}(); err != nil {
		return err
	}

	return nil
}
