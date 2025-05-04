package changedocumentsscreen

import (
	"context"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/screen"
	baseScreen "github.com/mkrtychanr/rag_bot/internal/screen/base"
	selectorscreen "github.com/mkrtychanr/rag_bot/internal/screen/selector/base"
	endtoend "github.com/mkrtychanr/rag_bot/internal/screen/selector/end_to_end"
)

var _ screen.Screen = &ChangeDocumentsScreen{}

type ChangeDocumentsScreen struct {
	endtoend.End2EndSelector
	getter getter
}

func NewChangeDocumentsScreen(g getter, base baseScreen.Base) *ChangeDocumentsScreen {
	obj := &ChangeDocumentsScreen{
		getter: g,
	}

	obj.End2EndSelector = endtoend.End2EndSelector{
		Selector: *selectorscreen.NewSelector(obj.wrapper, base),
	}

	return obj
}

func (s *ChangeDocumentsScreen) GetScreenType() screen.ScreenType {
	return screen.ChangeDocumentsScreen
}

func (s *ChangeDocumentsScreen) wrapper(ctx context.Context, id int64) ([]selectorscreen.SelectorData, error) {
	docs, err := s.getter.GetUserDocuments(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user documents: %w", err)
	}

	result := make([]selectorscreen.SelectorData, 0, len(docs))

	for _, doc := range docs {
		result = append(result, selectorscreen.SelectorData{
			ID:   doc.ID,
			Text: doc.Name,
		})
	}

	return result, nil
}
