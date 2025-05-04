package endtoend

import (
	"context"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/model"
	"github.com/mkrtychanr/rag_bot/internal/screen"
	baseSelector "github.com/mkrtychanr/rag_bot/internal/screen/selector/base"
)

type End2EndSelector struct {
	baseSelector.Selector
	performFunc func(context.Context, map[string]any) error
}

func NewE2ESelector(s baseSelector.Selector) *End2EndSelector {
	return &End2EndSelector{
		Selector: s,
	}
}

func NewEndlessSelector(s baseSelector.Selector, pf func(context.Context, map[string]any) error) *End2EndSelector {
	return &End2EndSelector{
		Selector:    s,
		performFunc: pf,
	}
}

func (s *End2EndSelector) Perform(ctx context.Context, payload map[string]any) (screen.Screen, error) {
	v, ok := payload["selector_option"].(baseSelector.SelectorOption)
	if !ok {
		return nil, screen.ErrWrongType
	}

	if v.Option == -1 || v.Option == -2 {
		payload["option"] = model.MenuOption{
			Option: int64(v.Option),
		}
		sc, err := s.Next(ctx, payload)
		if err != nil {
			return nil, fmt.Errorf("failed to next: %w", err)
		}

		return sc, nil
	}
	if s.performFunc != nil {
		if err := s.performFunc(ctx, payload); err != nil {
			return nil, fmt.Errorf("failed to run perform func: %w", err)
		}
	}

	nextScreen := s.NextScreens[0]

	if err := nextScreen.Load(ctx, payload); err != nil {
		return nil, fmt.Errorf("failed to load next screen: %w", err)
	}

	return nextScreen, nil
}

func (s *End2EndSelector) GetScreenType() screen.ScreenType {
	return screen.E2ESelectorScreen
}
