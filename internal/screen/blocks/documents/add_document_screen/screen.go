package adddocumentscreen

import (
	"context"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/model"
	"github.com/mkrtychanr/rag_bot/internal/screen"
	baseScreen "github.com/mkrtychanr/rag_bot/internal/screen/base"
)

var _ screen.Screen = &AddDocumentScreen{}

type AddDocumentScreen struct {
	baseScreen.Base
	documentAdder documentAdder
}

func NewAddDocumentScreen(da documentAdder, base baseScreen.Base) *AddDocumentScreen {
	return &AddDocumentScreen{
		documentAdder: da,
		Base:          base,
	}
}

func (s *AddDocumentScreen) Render() (model.Screen, error) {
	buttons, err := s.BuildBaseButtons()
	if err != nil {
		return model.Screen{}, fmt.Errorf("failed to build base buttons: %w", err)
	}

	return model.Screen{
		Text:    "Загрузите документ и напишите его название одним сообщением",
		Buttons: buttons,
	}, nil
}

func (s *AddDocumentScreen) GetTitle() string {
	return "Загрузить документ"
}

func (s *AddDocumentScreen) Load(_ context.Context, _ any) error {
	return nil
}

func (s *AddDocumentScreen) Perform(ctx context.Context, payload any) (screen.Screen, error) {
	if payload == nil {
		return nil, screen.ErrEmptyPayload
	}

	v, ok := payload.(PerformModel)
	if !ok {
		return nil, screen.ErrWrongType
	}

	if err := s.documentAdder.AddDocument(ctx, v.Name, v.UserID, v.DocumentID); err != nil {
		return nil, fmt.Errorf("failed to add document: %w", err)
	}

	return s.PreviousScreen, nil
}

func (s *AddDocumentScreen) GetScreenType() screen.ScreenType {
	return screen.AddDocumentScreen
}
