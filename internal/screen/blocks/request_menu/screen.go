package requestmenu

import (
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/model"
	baseScreen "github.com/mkrtychanr/rag_bot/internal/screen/base"
)

type RequestMenu struct {
	baseScreen.Base
	getData func() (string, error)
}

func (s *RequestMenu) Render() (model.Screen, error) {
	buttons, err := s.BuildBaseButtons()
	if err != nil {
		return model.Screen{}, fmt.Errorf("failed to build base buttons: %w", err)
	}

	text, err := s.getData()
	if err != nil {
		return model.Screen{}, fmt.Errorf("failed to get data: %w", err)
	}

	return model.Screen{
		Text:    text,
		Buttons: buttons,
	}, nil
}
