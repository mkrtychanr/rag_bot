package adddocumentintogroupscreen

import (
	"context"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/screen"
	baseScreen "github.com/mkrtychanr/rag_bot/internal/screen/base"
	selectorscreen "github.com/mkrtychanr/rag_bot/internal/screen/selector/base"
	endtoend "github.com/mkrtychanr/rag_bot/internal/screen/selector/end_to_end"
)

// type _ scr

type AddDocumentIntoGroupScreen struct {
	endtoend.End2EndSelector
	adder  adder
	getter getter
}

func NewAddDocumentIntoGroupScreen(a adder, g getter, base baseScreen.Base) *AddDocumentIntoGroupScreen {
	obj := &AddDocumentIntoGroupScreen{
		adder:  a,
		getter: g,
	}

	obj.End2EndSelector = *endtoend.NewEndlessSelector(
		*selectorscreen.NewSelector(obj.getterWrapper, base),
		obj.adderWrapper,
	)

	return obj
}

func (s *AddDocumentIntoGroupScreen) getterWrapper(ctx context.Context, id int64) ([]selectorscreen.SelectorData, error) {
	userID, ok := s.CurrentPayload["user_id"].(int64)
	if !ok {
		return nil, screen.ErrWrongType
	}

	paperID, ok := s.CurrentPayload["paper_id"].(int64)
	if !ok {
		return nil, screen.ErrWrongType
	}

	groups, err := s.getter.GetGroupsToAddPaper(ctx, userID, paperID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user documents: %w", err)
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

func (s *AddDocumentIntoGroupScreen) adderWrapper(ctx context.Context, payload map[string]any) error {
	paperID, ok := payload["paper_id"].(int64)
	if !ok {
		return screen.ErrWrongType
	}

	groupID, ok := payload["group_id"].(int64)
	if !ok {
		return screen.ErrWrongType
	}

	if err := s.adder.AddDocumentIntoGroup(ctx, paperID, groupID); err != nil {
		return fmt.Errorf("failed to add document into group: %w", err)
	}

	return nil
}

func (s *AddDocumentIntoGroupScreen) GetScreenType() screen.ScreenType {
	return screen.AddDocumentIntoGroupScreen
}
