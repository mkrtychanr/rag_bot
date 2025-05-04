package groupaccess

import (
	"context"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/model"
	"github.com/mkrtychanr/rag_bot/internal/utils"
)

type getGroupsFunc = func(context.Context, int64) ([]int64, error)

type useCase struct {
	repository repository
}

func NewUseCase(repository repository) *useCase {
	return &useCase{
		repository: repository,
	}
}

func (u *useCase) GetUserAccessGroups(ctx context.Context, tgID int64) ([]model.Group, error) {
	return u.getUserGroups(ctx, tgID, u.repository.GetUserGroups)
}

func (u *useCase) GetUserRWAccessGroups(ctx context.Context, tgID int64) ([]model.Group, error) {
	return u.getUserGroups(ctx, tgID, u.repository.GetUserRWGroupIDs)
}

func (u *useCase) getUserGroups(ctx context.Context, tgID int64, f getGroupsFunc) ([]model.Group, error) {
	groupIDs, err := f(ctx, tgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group ids: %w", err)
	}

	groups, err := u.repository.FetchGroupsInfo(ctx, groupIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch groups info: %w", err)
	}

	ownedGroups, err := u.repository.GetUserGroupsOwnership(ctx, tgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user groups ownership: %w", err)
	}

	m := make(map[int64]model.Group, len(groups)+len(ownedGroups))

	for _, group := range groups {
		m[group.ID] = group
	}

	for _, group := range ownedGroups {
		m[group.ID] = group
	}

	result := make([]model.Group, 0, len(m))

	for _, group := range m {
		result = append(result, group)
	}

	return result, nil
}

func (u *useCase) GetGroupsToAddPaper(ctx context.Context, tgID int64, paperID int64) ([]model.Group, error) {
	userAccessGroups, err := u.GetUserRWAccessGroups(ctx, tgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user rw access groups: %w", err)
	}

	paperGroupIDs, err := u.repository.GetPaperGroupIDs(ctx, paperID)
	if err != nil {
		return nil, fmt.Errorf("failed to get paper group ids: %w", err)
	}

	restrictedIDs := utils.SliceToMap(paperGroupIDs)

	groups := make([]model.Group, 0, len(userAccessGroups))

	for _, group := range userAccessGroups {
		if _, ok := restrictedIDs[group.ID]; ok {
			continue
		}

		groups = append(groups, group)
	}

	return groups, nil
}

func (u *useCase) GetGroupsToDeletePaper(ctx context.Context, tgID int64, paperID int64) ([]model.Group, error) {
	userAccessGroups, err := u.GetUserRWAccessGroups(ctx, tgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user rw access groups: %w", err)
	}

	paperGroupIDs, err := u.repository.GetPaperGroupIDs(ctx, paperID)
	if err != nil {
		return nil, fmt.Errorf("failed to get paper group ids: %w", err)
	}

	unrestrictedIDs := utils.SliceToMap(paperGroupIDs)

	groups := make([]model.Group, 0, len(userAccessGroups))

	for _, group := range userAccessGroups {
		if _, ok := unrestrictedIDs[group.ID]; !ok {
			continue
		}

		groups = append(groups, group)
	}

	return groups, nil
}
