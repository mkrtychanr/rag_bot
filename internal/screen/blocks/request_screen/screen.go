package requestscreen

import (
	"context"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/model"
	"github.com/mkrtychanr/rag_bot/internal/screen"
	baseScreen "github.com/mkrtychanr/rag_bot/internal/screen/base"
)

var _ screen.Screen = &RequestScreen{}

type RequestScreen struct {
	baseScreen.Base
	requestMaker requestMaker
	sender       sender
}

func NewRequestScreen(rm requestMaker, s sender, base baseScreen.Base) *RequestScreen {
	return &RequestScreen{
		requestMaker: rm,
		sender:       s,
		Base:         base,
	}
}

func (s *RequestScreen) Render() (model.Screen, error) {
	buttons, err := s.BuildBaseButtons()
	if err != nil {
		return model.Screen{}, fmt.Errorf("failed to build base buttons: %w", err)
	}

	return model.Screen{
		Text:    "Введите ваш запрос",
		Buttons: buttons,
	}, nil
}

func (s *RequestScreen) GetTitle() string {
	return "Сделать запрос"
}

func (s *RequestScreen) Load(ctx context.Context, payload map[string]any) error {
	return nil
}

func (s *RequestScreen) Perform(ctx context.Context, payload map[string]any) (screen.Screen, error) {
	if payload == nil {
		return nil, screen.ErrEmptyPayload
	}

	v, ok := payload["perform"].(PerformModel)
	if !ok {
		return nil, screen.ErrWrongType
	}

	result, err := s.requestMaker.MakeRequest(ctx, v.Request, v.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if err := s.sender.SendText(ctx, v.UserID, result); err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return s, nil
}

func (s *RequestScreen) GetScreenType() screen.ScreenType {
	return screen.RequestScreen
}
