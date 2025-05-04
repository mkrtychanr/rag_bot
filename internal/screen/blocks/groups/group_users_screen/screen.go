package groupusersscreen

import (
	"context"
	"fmt"
	"strings"

	"github.com/mkrtychanr/rag_bot/internal/model"
	"github.com/mkrtychanr/rag_bot/internal/screen"
	baseScreen "github.com/mkrtychanr/rag_bot/internal/screen/base"
	requestmenu "github.com/mkrtychanr/rag_bot/internal/screen/blocks/request_menu"
)

type GroupUsersScreen struct {
	requestmenu.RequestMenu
	getter getter
}

func NewGroupUsersScreen(g getter, base baseScreen.Base) *GroupUsersScreen {
	obj := &GroupUsersScreen{
		RequestMenu: requestmenu.RequestMenu{
			Base: base,
		},
		getter: g,
	}

	obj.GetData = obj.getDataWrapper

	return obj
}

func (s *GroupUsersScreen) getDataWrapper(ctx context.Context, payload map[string]any) (string, error) {
	groupID, ok := payload["group_id"].(int64)
	if !ok {
		return "", screen.ErrWrongType
	}

	users, err := s.getter.GetGroupUsers(ctx, groupID)
	if err != nil {
		return "", fmt.Errorf("failed to get group users: %w", err)
	}

	b := strings.Builder{}

	b.WriteString("Имя пользователя – уровень доступа\n")

	for _, user := range users {
		b.WriteString(fmt.Sprintf("@%s – %s\n", user.Shortname, getAccessTypeString(user.AccessType)))
	}

	return b.String(), nil
}

func getAccessTypeString(accessType model.RightsPolicy) string {
	switch accessType {
	case model.ReadOnlyRightPolicy:
		return "Только чтение"
	case model.ReadWriteRightPolicy:
		return "Чтение и запись"
	case model.Onwer:
		return "Хозяин группы"
	default:
		return "Неизвестно"
	}
}

func (s *GroupUsersScreen) GetScreenType() screen.ScreenType {
	return screen.GroupUsersScreen
}
