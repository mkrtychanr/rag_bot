package mainmenu

import (
	"context"

	"github.com/mkrtychanr/rag_bot/internal/screen"
	baseScreen "github.com/mkrtychanr/rag_bot/internal/screen/base"
)

var _ screen.Screen = &DefaultMenuScreen{}

type DefaultMenuScreen struct {
	baseScreen.Base
}

func (s *DefaultMenuScreen) Load(_ context.Context, _ map[string]any) error {
	return nil
}

func (s *DefaultMenuScreen) Perform(_ context.Context, _ map[string]any) (screen.Screen, error) {
	return nil, nil
}

func (s *DefaultMenuScreen) GetScreenType() screen.ScreenType {
	return screen.DefaultScreen
}
