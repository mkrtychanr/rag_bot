package groupaccessscreen

import (
	"context"
	"fmt"
	"strings"

	"github.com/mkrtychanr/rag_bot/internal/screen"
	baseScreen "github.com/mkrtychanr/rag_bot/internal/screen/base"
	requestmenu "github.com/mkrtychanr/rag_bot/internal/screen/blocks/request_menu"
)

type GroupAccessScreen struct {
	requestmenu.RequestMenu
	getter getter
}

func NewGroupAccessScreen(g getter, base baseScreen.Base) *GroupAccessScreen {
	obj := &GroupAccessScreen{
		RequestMenu: requestmenu.RequestMenu{
			Base: base,
		},
		getter: g,
	}

	obj.GetData = obj.getDataWrapper

	return obj
}

func (s *GroupAccessScreen) getDataWrapper(ctx context.Context, payload map[string]any) (string, error) {
	userID, ok := payload["user_id"].(int64)
	if !ok {
		return "", screen.ErrWrongType
	}

	groups, err := s.getter.GetUserAccessGroups(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user access groups: %w", err)
	}

	b := strings.Builder{}

	b.WriteString("Название – владелец\n")

	for _, group := range groups {
		b.WriteString(fmt.Sprintf("%s – @%s\n", group.Name, group.Admin.Shortname))
	}

	return b.String(), nil
}

func (s *GroupAccessScreen) GetScreenType() screen.ScreenType {
	return screen.GroupAccessScreen
}
