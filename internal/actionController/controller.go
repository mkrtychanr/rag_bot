package actioncontroller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/model"
	"github.com/mkrtychanr/rag_bot/internal/screen"
)

var ErrUnknownScreenType = errors.New("unknown screen type")

type controllerFunc = func(ctx context.Context, screen screen.Screen, data model.OperatorData) (screen.Screen, error)

type actionController struct {
	controllers map[screen.ScreenType]controllerFunc
}

func NewActionController() *actionController {
	return &actionController{
		controllers: map[screen.ScreenType]controllerFunc{
			screen.DefaultScreen: defaultScreenController,
			screen.RequestScreen: defaultScreenController,
		},
	}
}

func (ac *actionController) GetScreenController(screenType screen.ScreenType) (controllerFunc, error) {
	f, ok := ac.controllers[screenType]
	if !ok {
		return nil, ErrUnknownScreenType
	}

	return f, nil
}

func defaultScreenController(ctx context.Context, sc screen.Screen, payload model.OperatorData) (screen.Screen, error) {
	var v model.MenuOption
	if err := json.Unmarshal(payload.CallbackData, &v); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	newSc, err := sc.Next(ctx, v)
	if err != nil {
		return nil, fmt.Errorf("failed to get next screen: %w", err)
	}

	return newSc, nil
}
