package pickgroupscreen

import (
	"context"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/screen"
	baseScreen "github.com/mkrtychanr/rag_bot/internal/screen/base"
	selectorscreen "github.com/mkrtychanr/rag_bot/internal/screen/selector/base"
	endtoend "github.com/mkrtychanr/rag_bot/internal/screen/selector/end_to_end"
)

type PickGroupScreen struct {
	endtoend.End2EndSelector
	getter    getter
	performer performer
}

func NewPickGroupScreen(g getter, p performer, base baseScreen.Base) *PickGroupScreen {
	obj := &PickGroupScreen{
		getter:    g,
		performer: p,
	}

	s := *selectorscreen.NewSelector(obj.getterWrapper, base)

	selector := *endtoend.NewE2ESelector(s)

	if p != nil {
		selector = *endtoend.NewEndlessSelector(s, p.Perform)
	}

	obj.End2EndSelector = selector

	return obj
}

func (s *PickGroupScreen) getterWrapper(ctx context.Context, id int64) ([]selectorscreen.SelectorData, error) {
	userID, ok := s.CurrentPayload["user_id"].(int64)
	if !ok {
		return nil, screen.ErrWrongType
	}

	groups, err := s.getter.GetGroups(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get groups: %w", err)
	}

	result := make([]selectorscreen.SelectorData, 0, len(groups))

	for _, group := range groups {
		result = append(result, selectorscreen.SelectorData{
			ID:   group.ID,
			Text: group.Name,
		})
	}

	return result, nil
}

func (s *PickGroupScreen) GetScreenType() screen.ScreenType {
	return screen.PickGroupScreen
}
