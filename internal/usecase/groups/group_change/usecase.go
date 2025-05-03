package groupchange

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

func (u *useCase) ChangeGroupName(ctx context.Context, groupID int64, name string) error {
	if err := u.repository.ChangeGroupName(ctx, groupID, name); err != nil {
		return fmt.Errorf("failed to change group name: %w", err)
	}

	return nil
}

func (u *useCase) AddUserIntoGroup(ctx context.Context, groupID int64, tgID int64) error {
	if err := u.repository.AddUserIntoGroup(ctx, groupID, tgID); err != nil {
		return fmt.Errorf("failed to add user into group: %w", err)
	}

	return nil
}

func (u *useCase) getUsersWithRightPolicy(ctx context.Context, getUsersFunc func(ctx context.Context, groupID int64) ([]int64, error), groupID int64) ([]model.User, error) {
	ids, err := getUsersFunc(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user ids: %w", err)
	}

	users, err := u.repository.FetchUsersInfo(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users info: %w", err)
	}

	return users, nil
}

func (u *useCase) GetUsersWithReadOnlyRightPolicy(ctx context.Context, groupID int64) ([]model.User, error) {
	users, err := u.getUsersWithRightPolicy(ctx, u.repository.GetUserIDsInGroupWithReadOnlyAccess, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get users with right policy: %w", err)
	}

	return users, nil
}

func (u *useCase) GetUsersWithReadWriteRightPolicy(ctx context.Context, groupID int64) ([]model.User, error) {
	users, err := u.getUsersWithRightPolicy(ctx, u.repository.GetUserIDsInGroupWithReadWriteAccess, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get users with right policy: %w", err)
	}

	return users, nil
}

func (u *useCase) DeleteUserFromGroup(ctx context.Context, groupID int64, tgID int64) error {
	if err := u.repository.DeleteUserFromGroup(ctx, groupID, tgID); err != nil {
		return fmt.Errorf("failed to delete user from group: %w", err)
	}

	return nil
}
