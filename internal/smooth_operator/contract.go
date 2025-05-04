package smoothoperator

import (
	"context"

	"github.com/mkrtychanr/rag_bot/internal/model"
	"github.com/mkrtychanr/rag_bot/internal/screen"
)

type screenChanger interface {
	ChangeScreen(ctx context.Context, chatID int64, tgID int64, newScreen model.Screen) error
	Clear(ctx context.Context, chatID int64) error
}

type controllerFunc = func(ctx context.Context, screen screen.Screen, data model.OperatorData) (screen.Screen, error)

type actionController interface {
	GetScreenController(screenType screen.ScreenType) (controllerFunc, error)
}
