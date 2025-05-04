package groupinfo

import (
	"context"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/model"
)

type useCase struct {
	repository repository
}

func NewUseCase(repository repository) *useCase {
	return &useCase{
		repository: repository,
	}
}

func (u *useCase) GetGroupDocuments(ctx context.Context, groupID int64) ([]model.Paper, error) {
	paperIDs, err := u.repository.GetGroupPapers(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group papers: %w", err)
	}

	papers, err := u.repository.FetchPapersInfo(ctx, paperIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch papers info: %w", err)
	}

	return papers, nil
}

func (u *useCase) GetGroupUsers(ctx context.Context, groupID int64) ([]model.UserGroup, error) {
	users, err := u.repository.GetGroupUsers(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group users: %w", err)
	}

	groupInfo, err := u.repository.FetchGroupsInfo(ctx, []int64{groupID})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch groups info: %w", err)
	}

	owner, err := u.repository.FetchUsersInfo(ctx, []int64{groupInfo[0].Admin.ID})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users info: %w", err)
	}

	users = append(users, model.UserGroup{
		User:       owner[0],
		AccessType: model.Onwer,
	})

	return users, nil
}
