package deletedocumentscreen

import (
	"context"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/screen"
	baseScreen "github.com/mkrtychanr/rag_bot/internal/screen/base"
	selectorscreen "github.com/mkrtychanr/rag_bot/internal/screen/selector/base"
	endtoend "github.com/mkrtychanr/rag_bot/internal/screen/selector/end_to_end"
)

type DeleteDocumentScreen struct {
	endtoend.End2EndSelector
	deleter deleter
	getter  getter
}

func NewDeleteDocumentScreen(d deleter, g getter, base baseScreen.Base) *DeleteDocumentScreen {
	obj := &DeleteDocumentScreen{
		deleter: d,
		getter:  g,
	}

	obj.End2EndSelector = *endtoend.NewEndlessSelector(
		*selectorscreen.NewSelector(obj.getterWrapper, base),
		obj.deleterWrapper,
	)

	return obj
}

func (s *DeleteDocumentScreen) getterWrapper(ctx context.Context, id int64) ([]selectorscreen.SelectorData, error) {
	userID, ok := s.CurrentPayload["user_id"].(int64)
	if !ok {
		return nil, screen.ErrWrongType
	}

	docs, err := s.getter.GetUserDocuments(ctx, userID)
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

func (s *DeleteDocumentScreen) deleterWrapper(ctx context.Context, payload map[string]any) error {
	paperID, ok := payload["paper_id"].(int64)
	if !ok {
		return screen.ErrWrongType
	}

	if err := s.deleter.DeleteDocument(ctx, paperID); err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	return nil
}

func (s *DeleteDocumentScreen) GetScreenType() screen.ScreenType {
	return screen.DeleteDocumentScreen
}
