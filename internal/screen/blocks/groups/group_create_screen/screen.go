package groupcreatescreen

import (
	"context"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/model"
	"github.com/mkrtychanr/rag_bot/internal/screen"
	baseScreen "github.com/mkrtychanr/rag_bot/internal/screen/base"
)

type CreateGroupScreen struct {
	baseScreen.Base
	creator creator
}

func NewCreateGroupScreen(c creator, base baseScreen.Base) *CreateGroupScreen {
	return &CreateGroupScreen{
		creator: c,
		Base:    base,
	}
}

func (s *CreateGroupScreen) Render() (model.Screen, error) {
	buttons, err := s.BuildBaseButtons()
	if err != nil {
		return model.Screen{}, fmt.Errorf("failed to build base buttons: %w", err)
	}

	return model.Screen{
		Text:    "Введите название новой группы",
		Buttons: buttons,
	}, nil
}

func (s *CreateGroupScreen) GetTitle() string {
	return "Создать новую группу"
}

func (s *CreateGroupScreen) Load(_ context.Context, payload map[string]any) error {
	s.CurrentPayload = payload

	return nil
}

func (s *CreateGroupScreen) Perform(ctx context.Context, payload map[string]any) (screen.Screen, error) {
	if payload == nil {
		return nil, screen.ErrEmptyPayload
	}

	groupName, ok := payload["text"].(string)
	if !ok {
		return nil, screen.ErrWrongType
	}

	userID, ok := payload["user_id"].(int64)
	if !ok {
		return nil, screen.ErrWrongType
	}

	if err := s.creator.CreateGroup(ctx, userID, groupName); err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	payload["option"] = model.MenuOption{
		Option: -1,
	}

	return s.Next(ctx, payload)
}

func (s *CreateGroupScreen) GetScreenType() screen.ScreenType {
	return screen.CreateGroupScreen
}
