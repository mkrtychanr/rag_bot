package base

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/mkrtychanr/rag_bot/internal/model"
	"github.com/mkrtychanr/rag_bot/internal/screen"
)

type Base struct {
	self           screen.Screen
	Text           string
	Title          string
	NextScreens    []screen.Screen
	CurrentPayload []byte
	HeadScreen     screen.Screen `json:"-"`
	PreviousScreen screen.Screen `json:"-"`
}

func (s *Base) Next(ctx context.Context, payload model.MenuOption) (screen.Screen, error) {
	if payload.Option == -2 {
		if s.HeadScreen == nil {
			return s.self, nil
		}

		return s.HeadScreen, nil
	}

	if payload.Option == -1 {
		prev := s.PreviousScreen
		if prev == nil {
			return s.self, nil
		}

		if err := prev.Load(ctx, nil); err != nil {
			return nil, fmt.Errorf("failed to load previous screen: %w", err)
		}

		return prev, nil
	}

	if payload.Option >= int64(len(s.NextScreens)) {
		return nil, screen.ErrUnknownScreen
	}

	next := s.NextScreens[payload.Option]
	if err := next.Load(ctx, payload.Payload); err != nil {
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

func (s *Base) Render() (model.Screen, error) {
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
			return model.Screen{}, fmt.Errorf("failed to marshal menu option: %w", err)
		}

		button.Payload = payload

		buttons = append(buttons, []model.Button{button})
	}

	baseButtons, err := s.BuildBaseButtons()
	if err != nil {
		return model.Screen{}, fmt.Errorf("failed to build base buttons: %w", err)
	}

	buttons = append(buttons, baseButtons...)

	return model.Screen{
		Text:    s.Text,
		Buttons: buttons,
	}, nil
}

func (s *Base) GetTitle() string {
	return s.Title
}
