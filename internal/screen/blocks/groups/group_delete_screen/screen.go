package groupdeletescreen

import (
	"context"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/screen"
	baseScreen "github.com/mkrtychanr/rag_bot/internal/screen/base"
	pickgroupscreen "github.com/mkrtychanr/rag_bot/internal/screen/blocks/groups/pick_group_screen"
)

type DeleteGroupScreen struct {
	pickgroupscreen.PickGroupScreen
	deleter deleter
}

func NewDeleteGroupScreen(g getter, d deleter, base baseScreen.Base) *DeleteGroupScreen {
	obj := &DeleteGroupScreen{
		deleter: d,
	}

	obj.PickGroupScreen = *pickgroupscreen.NewPickGroupScreen(g, obj.deleterWrapper, base)

	return obj
}

func (s *DeleteGroupScreen) deleterWrapper(ctx context.Context, payload map[string]any) error {
	groupID, ok := payload["group_id"].(int64)
	if !ok {
		return screen.ErrWrongType
	}

	if err := s.deleter.DeleteGroup(ctx, groupID); err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}

	return nil
}
