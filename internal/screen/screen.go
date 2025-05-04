package screen

import (
	"context"

	"github.com/mkrtychanr/rag_bot/internal/model"
)

type Screen interface {
	Next(ctx context.Context, payload map[string]any) (Screen, error) // In Base
	GetTitle() string                                                 // In Base
	Load(ctx context.Context, payload map[string]any) error           // In Base Selector + Base AfterSelector
	Render() (model.Screen, error)                                    // In Base + Base Selector + Base AfterSelector + Unique
	ExtractPayload() map[string]any
	Perform(ctx context.Context, payload map[string]any) (Screen, error)
	GetScreenType() ScreenType // unique only
}
