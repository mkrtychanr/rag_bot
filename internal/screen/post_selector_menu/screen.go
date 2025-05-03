package postselectormenu

import (
	"encoding/json"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/model"
	"github.com/mkrtychanr/rag_bot/internal/screen"
	baseScreen "github.com/mkrtychanr/rag_bot/internal/screen/base"
)

type PostSelectorMenu struct {
	baseScreen.Base
	Text string
}

func (s *PostSelectorMenu) Load(payload []byte) error {
	if payload == nil {
		return screen.ErrEmptyPayload
	}

	var v LoadModel
	if err := json.Unmarshal(payload, &v); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	s.Text = v.Text

	return nil
}

func (s *PostSelectorMenu) Render() (model.Screen, error) {
	buttons, err := s.BuildBaseButtons()
	if err != nil {
		return model.Screen{}, fmt.Errorf("failed to build base buttons: %w", err)
	}

	return model.Screen{
		Text:    "Выберите действие над " + s.Text,
		Buttons: buttons,
	}, nil
}
