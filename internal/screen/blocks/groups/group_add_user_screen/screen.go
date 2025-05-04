package groupadduserscreen

import (
	"context"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/model"
	"github.com/mkrtychanr/rag_bot/internal/screen"
	baseScreen "github.com/mkrtychanr/rag_bot/internal/screen/base"
)

type AddUserIntoGroupScreen struct {
	baseScreen.Base
	adder adder
}

func NewAddUserIntoGroupScreen(a adder, base baseScreen.Base) *AddUserIntoGroupScreen {
	return &AddUserIntoGroupScreen{
		adder: a,
		Base:  base,
	}
}

func (s *AddUserIntoGroupScreen) Render() (model.Screen, error) {
	buttons, err := s.BuildBaseButtons()
	if err != nil {
		return model.Screen{}, fmt.Errorf("failed to build base buttons: %w", err)
	}

	return model.Screen{
		Text:    "Введите имя пользователя",
		Buttons: buttons,
	}, nil
}

func (s *AddUserIntoGroupScreen) GetTitle() string {
	return "Добавить пользователя"
}

func (s *AddUserIntoGroupScreen) Load(_ context.Context, payload map[string]any) error {
	s.CurrentPayload = payload

	return nil
}

func (s *AddUserIntoGroupScreen) Perform(ctx context.Context, payload map[string]any) (screen.Screen, error) {
	if payload == nil {
		return nil, screen.ErrEmptyPayload
	}

	userName, ok := payload["text"].(string)
	if !ok {
		return nil, screen.ErrWrongType
	}

	groupID, ok := payload["group_id"].(int64)
	if !ok {
		return nil, screen.ErrWrongType
	}

	if err := s.adder.AddUserIntoGroup(ctx, groupID, userName); err != nil {
		return nil, fmt.Errorf("failed to change group name: %w", err)
	}

	payload["option"] = model.MenuOption{
		Option: -1,
	}

	return s.Next(ctx, payload)
}

func (s *AddUserIntoGroupScreen) GetScreenType() screen.ScreenType {
	return screen.AddUserIntoGroupScreen
}
