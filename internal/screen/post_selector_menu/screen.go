package postselectormenu

import (
	"context"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/model"
	"github.com/mkrtychanr/rag_bot/internal/screen"
	baseScreen "github.com/mkrtychanr/rag_bot/internal/screen/base"
	baseSelector "github.com/mkrtychanr/rag_bot/internal/screen/selector/base"
)

type PostSelectorMenu struct {
	baseScreen.Base
	Text string
}

func (s *PostSelectorMenu) Load(_ context.Context, payload map[string]any) error {
	s.CurrentPayload = payload
	if payload == nil {
		return screen.ErrEmptyPayload
	}

	v, ok := payload["selector_option"].(baseSelector.SelectorOption)
	if !ok {
		return screen.ErrWrongType
	}

	s.Text = v.Payload.Text

	return nil
}

func (s *PostSelectorMenu) Render() (model.Screen, error) {
	bs, err := s.Base.Render()
	if err != nil {
		return model.Screen{}, fmt.Errorf("failed to render base: %w", err)
	}

	bs.Text = "Выберите действие над " + s.Text

	return bs, nil
}

func (s *PostSelectorMenu) Perform(_ context.Context, _ map[string]any) (screen.Screen, error) {
	return nil, nil
}

func (s *PostSelectorMenu) GetScreenType() screen.ScreenType {
	return screen.PostE2ESelectorScreen
}
