package base

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/mkrtychanr/rag_bot/internal/model"
	"github.com/mkrtychanr/rag_bot/internal/screen"
	"github.com/mkrtychanr/rag_bot/internal/utils"
)

type Base struct {
	Text           string
	Title          string
	NextScreens    []screen.Screen
	CurrentPayload map[string]any
	HeadScreen     screen.Screen `json:"-"`
	PreviousScreen screen.Screen `json:"-"`
}

func (s *Base) Next(ctx context.Context, payload map[string]any) (screen.Screen, error) {
	option, ok := payload["option"].(model.MenuOption)
	if !ok {
		return nil, screen.ErrWrongType
	}

	if option.Option == -2 {

		return s.HeadScreen, nil
	}

	if option.Option == -1 {
		prev := s.PreviousScreen

		if err := prev.Load(ctx, s.ExtractPayload()); err != nil {
			return nil, fmt.Errorf("failed to load previous screen: %w", err)
		}

		return prev, nil
	}

	if option.Option >= int64(len(s.NextScreens)) {
		return nil, screen.ErrUnknownScreen
	}

	next := s.NextScreens[option.Option]
	if err := next.Load(ctx, payload); err != nil {
		return nil, fmt.Errorf("failed to load next screen: %w", err)
	}

	return next, nil
}

func (s *Base) BuildBaseButtons() ([][]model.Button, error) {
	buttons := make([][]model.Button, 0, 2)

	if s.PreviousScreen != nil {
		op := model.MenuOption{
			Option: -1,
		}

		payload, err := json.Marshal(op)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal back button: %w", err)
		}

		buttons = append(buttons, []model.Button{
			{
				Text:    "Назад",
				Payload: payload,
			},
		})
	}

	if s.HeadScreen != nil && s.HeadScreen != s.PreviousScreen {
		op := model.MenuOption{
			Option: -2,
		}

		payload, err := json.Marshal(op)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal back button: %w", err)
		}

		buttons = append(buttons, []model.Button{
			{
				Text:    "В главное меню",
				Payload: payload,
			},
		})
	}

	return buttons, nil
}

func (s *Base) BuildNextScreensButtons() ([][]model.Button, error) {
	buttons := make([][]model.Button, 0, len(s.NextScreens))

	for i, screen := range s.NextScreens {
		button := model.Button{
			Text: screen.GetTitle(),
		}

		op := model.MenuOption{
			Option: int64(i),
		}

		payload, err := json.Marshal(op)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal menu option: %w", err)
		}

		button.Payload = payload

		buttons = append(buttons, []model.Button{button})
	}

	return buttons, nil
}

func (s *Base) Render() (model.Screen, error) {
	nextScreensButtons, err := s.BuildNextScreensButtons()
	if err != nil {
		return model.Screen{}, fmt.Errorf("failed to build next screens buttons: %w", err)
	}

	baseButtons, err := s.BuildBaseButtons()
	if err != nil {
		return model.Screen{}, fmt.Errorf("failed to build base buttons: %w", err)
	}

	buttons := append(nextScreensButtons, baseButtons...)

	return model.Screen{
		Text:    s.Text,
		Buttons: buttons,
	}, nil
}

func (s *Base) GetTitle() string {
	return s.Title
}

func (s *Base) ExtractPayload() map[string]any {
	return utils.Copy(s.CurrentPayload)
}
