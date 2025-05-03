package selectorscreen

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/mkrtychanr/rag_bot/internal/model"
	"github.com/mkrtychanr/rag_bot/internal/screen"
	baseScreen "github.com/mkrtychanr/rag_bot/internal/screen/base"
	"github.com/mkrtychanr/rag_bot/internal/utils"
)

const buttonsInARow = 5

type SelectorData struct {
	ID   int64
	Text string
}

type Selector struct {
	baseScreen.Base
	GetSelectorData func(context.Context, int64) ([]SelectorData, error)
	data            []SelectorData
}

func NewSelector(f func(context.Context, int64) ([]SelectorData, error), base baseScreen.Base) *Selector {
	return &Selector{
		GetSelectorData: f,
		Base:            base,
	}
}

func (s *Selector) Load(ctx context.Context, payload any) error {
	if payload == nil {
		return screen.ErrEmptyPayload
	}

	v, ok := payload.(LoadModel)
	if !ok {
		return screen.ErrWrongType
	}

	data, err := s.GetSelectorData(ctx, v.ID)
	if err != nil {
		return fmt.Errorf("failed to get selector data: %w", err)
	}

	s.data = data

	return nil
}

func (s *Selector) buildButtons() ([][]model.Button, error) {
	m := utils.NewMatrix[model.Button](buttonsInARow)

	for i := 0; i < len(s.data); i++ {
		b, err := json.Marshal(s.data[i])
		if err != nil {
			return nil, fmt.Errorf("failed to marshal selector data: %w", err)
		}

		button := model.Button{
			Text:    strconv.Itoa(i + 1),
			Payload: b,
		}

		m.Add(button)
	}

	baseButtons, err := s.BuildBaseButtons()
	if err != nil {
		return nil, fmt.Errorf("failed to build base buttons: %w", err)
	}

	buttons := m.GetMatrix()

	buttons = append(buttons, baseButtons...)

	return buttons, nil
}

func (s *Selector) buildText() (string, error) {
	b := strings.Builder{}

	if _, err := b.WriteString("Выберите вариант из списка:"); err != nil {
		return "", fmt.Errorf("failed to write header string: %w", err)
	}

	for i, v := range s.data {
		if _, err := b.WriteString(fmt.Sprintf("%d – %s", i+1, v.Text)); err != nil {
			return "", fmt.Errorf("failed to write block string: %w", err)
		}
	}

	return b.String(), nil
}

func (s *Selector) Render() (model.Screen, error) {
	text, err := s.buildText()
	if err != nil {
		return model.Screen{}, fmt.Errorf("failed to build text: %w", err)
	}

	buttons, err := s.buildButtons()
	if err != nil {
		return model.Screen{}, fmt.Errorf("failed to build buttons: %w", err)
	}

	return model.Screen{
		Text:    text,
		Buttons: buttons,
	}, nil
}
