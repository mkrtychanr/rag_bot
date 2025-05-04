package groupdocumentsscreen

import (
	"context"
	"fmt"
	"strings"

	"github.com/mkrtychanr/rag_bot/internal/screen"
	baseScreen "github.com/mkrtychanr/rag_bot/internal/screen/base"
	requestmenu "github.com/mkrtychanr/rag_bot/internal/screen/blocks/request_menu"
)

type GroupDocumentsScreen struct {
	requestmenu.RequestMenu
	getter getter
}

func NewGroupDocumentsScreen(g getter, base baseScreen.Base) *GroupDocumentsScreen {
	obj := &GroupDocumentsScreen{
		RequestMenu: requestmenu.RequestMenu{
			Base: base,
		},
		getter: g,
	}

	obj.GetData = obj.getDataWrapper

	return obj
}

func (s *GroupDocumentsScreen) getDataWrapper(ctx context.Context, payload map[string]any) (string, error) {
	groupID, ok := payload["group_id"].(int64)
	if !ok {
		return "", screen.ErrWrongType
	}

	docs, err := s.getter.GetGroupDocuments(ctx, groupID)
	if err != nil {
		return "", fmt.Errorf("failed to get group documents: %w", err)
	}

	b := strings.Builder{}

	b.WriteString("Имя документа\n")

	for _, doc := range docs {
		b.WriteString(fmt.Sprintf("%s\n", doc.Name))
	}

	return b.String(), nil
}

func (s *GroupDocumentsScreen) GetScreenType() screen.ScreenType {
	return screen.GroupDocumentsScreen
}
