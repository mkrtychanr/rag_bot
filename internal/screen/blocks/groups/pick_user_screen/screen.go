package pickuserscreen

import (
	"context"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/model"
	"github.com/mkrtychanr/rag_bot/internal/screen"
	baseScreen "github.com/mkrtychanr/rag_bot/internal/screen/base"
	selectorscreen "github.com/mkrtychanr/rag_bot/internal/screen/selector/base"
	endtoend "github.com/mkrtychanr/rag_bot/internal/screen/selector/end_to_end"
)

type getterFunc = func(context.Context, int64) ([]model.User, error)
type performerFunc = func(context.Context, int64, int64) error

type PickUserScreen struct {
	endtoend.End2EndSelector
	getter    getterFunc
	performer performerFunc
}

func NewPickUserScreen(g getterFunc, p performerFunc, base baseScreen.Base) *PickUserScreen {
	obj := &PickUserScreen{
		getter:    g,
		performer: p,
	}

	obj.End2EndSelector = *endtoend.NewEndlessSelector(
		*selectorscreen.NewSelector(obj.getterWrapper, base),
		obj.performerWrapper,
	)

	return obj
}

func (s *PickUserScreen) getterWrapper(ctx context.Context, id int64) ([]selectorscreen.SelectorData, error) {
	groupID, ok := s.CurrentPayload["group_id"].(int64)
	if !ok {
		return nil, screen.ErrWrongType
	}

	users, err := s.getter(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	result := make([]selectorscreen.SelectorData, 0, len(users))

	for _, user := range users {
		result = append(result, selectorscreen.SelectorData{
			ID:   user.ID,
			Text: user.Shortname,
		})
	}

	return result, nil
}

func (s *PickUserScreen) performerWrapper(ctx context.Context, payload map[string]any) error {
	groupID, ok := payload["group_id"].(int64)
	if !ok {
		return screen.ErrWrongType
	}

	userID, ok := payload["action_user_id"].(int64)
	if !ok {
		return screen.ErrWrongType
	}

	return s.performer(ctx, groupID, userID)
}

func (s *PickUserScreen) GetScreenType() screen.ScreenType {
	return screen.PickUserScreen
}
