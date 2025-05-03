package endtoend

import (
	"context"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/screen"
	baseSelector "github.com/mkrtychanr/rag_bot/internal/screen/selector/base"
)

type End2EndSelector struct {
	baseSelector.Selector
}

func (s *End2EndSelector) Perform(ctx context.Context, payload any) (screen.Screen, error) {
	nextScreen := s.NextScreens[0]

	if err := nextScreen.Load(ctx, payload); err != nil {
		return nil, fmt.Errorf("failed to load next screen: %w", err)
	}

	return nextScreen, nil
}

func (s *End2EndSelector) GetScreenType() screen.ScreenType {
	return screen.E2ESelectorScreen
}
