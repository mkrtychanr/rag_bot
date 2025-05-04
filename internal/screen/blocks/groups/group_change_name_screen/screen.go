package groupchangenamescreen

import (
	"context"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/model"
	"github.com/mkrtychanr/rag_bot/internal/screen"
	baseScreen "github.com/mkrtychanr/rag_bot/internal/screen/base"
)

type ChangeGroupNameScreen struct {
	baseScreen.Base
	changer changer
}

func NewChangeGroupNameScreen(c changer, base baseScreen.Base) *ChangeGroupNameScreen {
	return &ChangeGroupNameScreen{
		changer: c,
		Base:    base,
	}
}

func (s *ChangeGroupNameScreen) Render() (model.Screen, error) {
	buttons, err := s.BuildBaseButtons()
	if err != nil {
		return model.Screen{}, fmt.Errorf("failed to build base buttons: %w", err)
	}

	return model.Screen{
		Text:    "Введите новое название группы",
		Buttons: buttons,
	}, nil
}

func (s *ChangeGroupNameScreen) GetTitle() string {
	return "Изменить название группы"
}

func (s *ChangeGroupNameScreen) Load(_ context.Context, payload map[string]any) error {
	s.CurrentPayload = payload

	return nil
}

func (s *ChangeGroupNameScreen) Perform(ctx context.Context, payload map[string]any) (screen.Screen, error) {
	if payload == nil {
		return nil, screen.ErrEmptyPayload
	}

	groupName, ok := payload["text"].(string)
	if !ok {
		return nil, screen.ErrWrongType
	}

	groupID, ok := payload["group_id"].(int64)
	if !ok {
		return nil, screen.ErrWrongType
	}

	if err := s.changer.ChangeGroupName(ctx, groupID, groupName); err != nil {
		return nil, fmt.Errorf("failed to change group name: %w", err)
	}

	payload["option"] = model.MenuOption{
		Option: -1,
	}

	return s.Next(ctx, payload)
}

func (s *ChangeGroupNameScreen) GetScreenType() screen.ScreenType {
	return screen.ChangeGroupNameScreen
}
