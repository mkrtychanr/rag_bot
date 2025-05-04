package requestmenu

import (
	"context"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/model"
	"github.com/mkrtychanr/rag_bot/internal/screen"
	baseScreen "github.com/mkrtychanr/rag_bot/internal/screen/base"
)

type RequestMenu struct {
	baseScreen.Base
	GetData func(context.Context, map[string]any) (string, error)
	data    string
}

func (s *RequestMenu) Render() (model.Screen, error) {
	buttons, err := s.BuildBaseButtons()
	if err != nil {
		return model.Screen{}, fmt.Errorf("failed to build base buttons: %w", err)
	}

	return model.Screen{
		Text:    s.data,
		Buttons: buttons,
	}, nil
}

func (s *RequestMenu) Load(ctx context.Context, payload map[string]any) error {
	s.CurrentPayload = payload

	data, err := s.GetData(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to get data: %w", err)
	}

	s.data = data

	return nil
}

func (s *RequestMenu) Perform(_ context.Context, _ map[string]any) (screen.Screen, error) {
	return nil, nil
}
