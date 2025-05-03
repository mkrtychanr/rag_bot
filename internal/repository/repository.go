package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/mkrtychanr/rag_bot/internal/logger"
	"github.com/mkrtychanr/rag_bot/internal/model"
	"github.com/mkrtychanr/rag_bot/internal/utils"
)

var (
	ErrNoRows           = errors.New("no rows returned")
	ErrUnexpectedResult = errors.New("unexpected result")
)

type repository struct {
	storage storage
}

func NewRepository(s storage) *repository {
	return &repository{
		storage: s,
	}
}

func (r *repository) RegistrateUser(ctx context.Context, name string, tgID int64) (int64, error) {
	result, err := r.storage.Query(ctx, registrateUserQuery, name, tgID)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	if !result.Next() {
		return 0, ErrNoRows
	}

	var userID int64
	if err := result.Scan(&userID); err != nil {
		return 0, fmt.Errorf("failed to scan userID: %w", err)
	}

	return userID, nil
}

func (r *repository) GetUserIDByTelegramID(ctx context.Context, tgID int64) (int64, error) {
	result, err := r.storage.Query(ctx, getUserIDByTelegramIDQuery, tgID)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	if !result.Next() {
		return 0, ErrNoRows
	}

	var userID int64
	if err := result.Scan(&userID); err != nil {
		return 0, fmt.Errorf("failed to scan userID: %w", err)
	}

	return userID, nil
}

func (r *repository) CreateGroup(ctx context.Context, tgID int64, name string) error {
	userID, err := r.GetUserIDByTelegramID(ctx, tgID)
	if err != nil {
		return fmt.Errorf("failed to get user id by telegram id: %w", err)
	}

	result, err := r.storage.Query(ctx, createGroupQuery, name, userID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	if !result.Next() {
		return ErrNoRows
	}

	var groupID int64
	if err := result.Scan(&groupID); err != nil {
		return fmt.Errorf("failed to scan groupID: %w", err)
	}

	if groupID == 0 {
		return ErrUnexpectedResult
	}

	return nil
}

func (r *repository) AddUserIntoGroup(ctx context.Context, groupID int64, tgID int64) error {
	userID, err := r.GetUserIDByTelegramID(ctx, tgID)
	if err != nil {
		return fmt.Errorf("failed to get user id by telegram id: %w", err)
	}

	result, err := r.storage.Query(ctx, addUserIntoGroupQuery, groupID, userID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	if !result.Next() {
		return ErrNoRows
	}

	var groupUserID int64
	if err := result.Scan(&groupUserID); err != nil {
		return fmt.Errorf("failed to scan groupUserID: %w", err)
	}

	if groupUserID == 0 {
		return ErrUnexpectedResult
	}

	return nil
}

func (r *repository) SetReadOnlyRightsForUserInGroup(ctx context.Context, groupID int64, tgID int64) error {
	if err := r.updateUserRightsPolicy(ctx, groupID, tgID, model.ReadOnlyRightPolicy); err != nil {
		return fmt.Errorf("failed to update rights policy: %w", err)
	}

	return nil
}

func (r *repository) SetReadWriteRightsForUserInGroup(ctx context.Context, groupID int64, tgID int64) error {
	if err := r.updateUserRightsPolicy(ctx, groupID, tgID, model.ReadWriteRightPolicy); err != nil {
		return fmt.Errorf("failed to update rights policy: %w", err)
	}

	return nil
}

func (r *repository) updateUserRightsPolicy(ctx context.Context, groupID int64, tgID int64, rightsPolicy model.RightsPolicy) error {
	userID, err := r.GetUserIDByTelegramID(ctx, tgID)
	if err != nil {
		return fmt.Errorf("failed to get user id by telegram id: %w", err)
	}

	tx, err := r.storage.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	var rollback bool

	defer func() {
		if rollback {
			if err := tx.Rollback(ctx); err != nil {
				logger.GetLogger().Err(err).Msg("failed to rollback transaction")
			}
		}
	}()

	tag, err := tx.Exec(ctx, changeRightsPolicyQuery, groupID, userID, rightsPolicy)
	if err != nil {
		rollback = true

		return fmt.Errorf("failed to execute query: %w", err)
	}

	if tag.RowsAffected() != 1 {
		rollback = true
		return ErrUnexpectedResult
	}

	if err := tx.Commit(ctx); err != nil {
		rollback = true

		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *repository) AddPaper(ctx context.Context, name string, tgID int64) (int64, error) {
	userID, err := r.GetUserIDByTelegramID(ctx, tgID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user id by telegram id: %w", err)
	}

	result, err := r.storage.Query(ctx, addPaperQuery, name, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	if !result.Next() {
		return 0, ErrNoRows
	}

	var paperID int64
	if err := result.Scan(&paperID); err != nil {
		return 0, fmt.Errorf("failed to scan paperID: %w", err)
	}

	if paperID == 0 {
		return 0, ErrUnexpectedResult
	}

	return paperID, nil
}

func (r *repository) AddPaperIntoGroup(ctx context.Context, paperID int64, groupID int64) error {
	result, err := r.storage.Query(ctx, addPaperIntoGroupQuery, paperID, groupID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	if !result.Next() {
		return ErrNoRows
	}

	var paperGroupID int64
	if err := result.Scan(&paperGroupID); err != nil {
		return fmt.Errorf("failed to scan paperGroupID: %w", err)
	}

	if paperGroupID == 0 {
		return ErrUnexpectedResult
	}

	return nil
}

func (r *repository) DeletePapers(ctx context.Context, papersIDs []int64) error {
	tx, err := r.storage.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	var rollback bool

	defer func() {
		if rollback {
			if err := tx.Rollback(ctx); err != nil {
				logger.GetLogger().Err(err).Msg("failed to rollback transaction")
			}
		}
	}()

	if _, err := tx.Exec(ctx, deletePapersFromPaperGroupQuery, papersIDs); err != nil {
		rollback = true

		return fmt.Errorf("failed to delete papers from paper group: %w", err)
	}

	if _, err := tx.Exec(ctx, deletePapersFromPapersQuery, papersIDs); err != nil {
		rollback = true

		return fmt.Errorf("failed to delete papers from papers: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		rollback = true

		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *repository) DeleteGroups(ctx context.Context, groupIDs []int64) error {
	tx, err := r.storage.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	var rollback bool

	defer func() {
		if rollback {
			if err := tx.Rollback(ctx); err != nil {
				logger.GetLogger().Err(err).Msg("failed to rollback transaction")
			}
		}
	}()

	if _, err := tx.Exec(ctx, deleteGroupsFromGroupUserQuery, groupIDs); err != nil {
		rollback = true

		return fmt.Errorf("failed to delete groups from group user: %w", err)
	}

	if _, err := tx.Exec(ctx, deleteGroupsFromPaperGroupQuery, groupIDs); err != nil {
		rollback = true

		return fmt.Errorf("failed to delete groups from paper group: %w", err)
	}

	if _, err := tx.Exec(ctx, deleteGroupsFromGroupsQuery, groupIDs); err != nil {
		rollback = true

		return fmt.Errorf("failed to delete groups from groups: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		rollback = true

		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *repository) GetUserGroups(ctx context.Context, tgID int64) ([]int64, error) {
	userID, err := r.GetUserIDByTelegramID(ctx, tgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user id by telegram id: %w", err)
	}

	result, err := r.storage.Query(ctx, getUserGroupsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getUserGroupsQuery: %w", err)
	}

	groups := make([]int64, 0)

	for result.Next() {
		var group int64
		if err := result.Scan(&group); err != nil {
			return nil, fmt.Errorf("failed to scan groups: %w", err)
		}

		groups = append(groups, group)
	}

	return groups, nil
}

func (r *repository) FetchGroupsInfo(ctx context.Context, groupIDs []int64) ([]model.Group, error) {
	result, err := r.storage.Query(ctx, getGroupsInfoQuery, groupIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get groups info query: %w", err)
	}

	groups := make([]model.Group, 0, len(groupIDs))

	for result.Next() {
		var (
			id      int64
			name    string
			adminID int64
		)

		if err := result.Scan(&id, &name, &adminID); err != nil {
			return nil, fmt.Errorf("failed to scan groups info: %w", err)
		}

		groups = append(groups, model.Group{
			ID:   id,
			Name: name,
			Admin: model.User{
				ID: adminID,
			},
		})
	}

	return groups, nil
}

func (r *repository) GetUserPapers(ctx context.Context, tgID int64) ([]model.Paper, error) {
	userID, err := r.GetUserIDByTelegramID(ctx, tgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user id by telegram id: %w", err)
	}

	result, err := r.storage.Query(ctx, getUserPapersQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get user papers query: %w", err)
	}

	papers := make([]model.Paper, 0)

	for result.Next() {
		var (
			id   int64
			name string
		)

		if err := result.Scan(&id, &name); err != nil {
			return nil, fmt.Errorf("failed to scan user papers: %w", err)
		}

		papers = append(papers, model.Paper{
			ID:   id,
			Name: name,
			Author: model.User{
				ID:         userID,
				TelegramID: tgID,
			},
		})
	}

	return papers, nil
}

func (r *repository) DeletePaperFromGroup(ctx context.Context, paperID int64, groupID int64) error {
	tx, err := r.storage.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	var rollback bool

	defer func() {
		if rollback {
			if err := tx.Rollback(ctx); err != nil {
				logger.GetLogger().Err(err).Msg("failed to rollback transaction")
			}
		}
	}()

	tag, err := tx.Exec(ctx, deletePaperFromGroupQuery, paperID, groupID)
	if err != nil {
		rollback = true

		return fmt.Errorf("failed to execute delete paper from group query: %w", err)
	}

	if tag.RowsAffected() != 1 {
		rollback = true

		return ErrUnexpectedResult
	}

	if err := tx.Commit(ctx); err != nil {
		rollback = true

		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *repository) GetUserGroupsOwnership(ctx context.Context, tgID int64) ([]model.Group, error) {
	userID, err := r.GetUserIDByTelegramID(ctx, tgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user id by telegram id: %w", err)
	}

	result, err := r.storage.Query(ctx, getUserGroupsOwnershipQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get user groups ownership query: %w", err)
	}

	groups := make([]model.Group, 0)

	for result.Next() {
		var (
			id   int64
			name string
		)

		if err := result.Scan(&id, &name); err != nil {
			return nil, fmt.Errorf("failed to scan user groups ownership: %w", err)
		}

		groups = append(groups, model.Group{
			ID:   id,
			Name: name,
			Admin: model.User{
				ID:         userID,
				TelegramID: tgID,
			},
		})
	}

	return groups, nil
}

func (r *repository) GetGroupPapers(ctx context.Context, groupID int64) ([]int64, error) {
	result, err := r.storage.Query(ctx, getGroupPapersQuery, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get group papers query: %w", err)
	}

	papers := make([]int64, 0)

	for result.Next() {
		var id int64

		if err := result.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan paper id: %w", err)
		}

		papers = append(papers, id)
	}

	return papers, nil
}

func (r *repository) FetchPapersInfo(ctx context.Context, paperIDs []int64) ([]model.Paper, error) {
	result, err := r.storage.Query(ctx, getPapersInfoQuery, paperIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get papers info query: %w", err)
	}

	papers := make([]model.Paper, 0, len(paperIDs))

	for result.Next() {
		var (
			id       int64
			name     string
			authorID int64
		)

		if err := result.Scan(&id, &name, &authorID); err != nil {
			return nil, fmt.Errorf("failed to scan paper info: %w", err)
		}

		papers = append(papers, model.Paper{
			ID:   id,
			Name: name,
			Author: model.User{
				ID: authorID,
			},
		})
	}

	return papers, nil
}

func (r *repository) GetGroupUsers(ctx context.Context, groupID int64) ([]model.UserGroup, error) {
	result, err := r.storage.Query(ctx, getGroupUsersQuery, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get group users query: %w", err)
	}

	users := make(map[int64]model.UserGroup)

	for result.Next() {
		var (
			id         int64
			accessType model.RightsPolicy
		)

		if err := result.Scan(&id, &accessType); err != nil {
			return nil, fmt.Errorf("failed to scan group user: %w", err)
		}

		userGroup := model.UserGroup{}
		userGroup.ID = id
		userGroup.AccessType = accessType

		users[id] = userGroup
	}

	userIDs := utils.MapKeysToSlice(users)

	fetchedUsers, err := r.FetchUsersInfo(ctx, userIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users info: %w", err)
	}

	for _, user := range fetchedUsers {
		v := users[user.ID]

		v.Shortname = user.Shortname
		v.TelegramID = user.TelegramID

		users[user.ID] = v
	}

	return utils.MapValuesToSlice(users), nil
}

func (r *repository) FetchUsersInfo(ctx context.Context, userIDs []int64) ([]model.User, error) {
	result, err := r.storage.Query(ctx, getUsersInfoQuery, userIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get users info query: %w", err)
	}

	users := make([]model.User, 0)

	for result.Next() {
		var (
			id   int64
			name string
			tgID int64
		)

		if err := result.Scan(&id, &name, &tgID); err != nil {
			return nil, fmt.Errorf("failed to get user info: %w", err)
		}

		users = append(users, model.User{
			ID:         id,
			Shortname:  name,
			TelegramID: tgID,
		})
	}

	return users, nil
}

func (r *repository) ChangeGroupName(ctx context.Context, groupID int64, name string) error {
	tx, err := r.storage.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	var rollback bool

	defer func() {
		if rollback {
			if err := tx.Rollback(ctx); err != nil {
				logger.GetLogger().Err(err).Msg("failed to rollback transaction")
			}
		}
	}()

	tag, err := tx.Exec(ctx, changeGroupNameQuery, groupID, name)
	if err != nil {
		rollback = true

		return fmt.Errorf("failed to execute change group name query: %w", err)
	}

	if tag.RowsAffected() != 1 {
		rollback = true

		return ErrUnexpectedResult
	}

	if err := tx.Commit(ctx); err != nil {
		rollback = true

		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *repository) GetUserByShortname(ctx context.Context, shortname string) (model.User, error) {
	result, err := r.storage.Query(ctx, getUserByShortnameQuery, shortname)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to execute get user by shortname query: %w", err)
	}

	if !result.Next() {
		return model.User{}, ErrNoRows
	}

	var (
		id   int64
		name string
		tgID int64
	)

	if err := result.Scan(&id, &name, &tgID); err != nil {
		return model.User{}, fmt.Errorf("failed to scan user info: %w", err)
	}

	return model.User{
		ID:         id,
		Shortname:  name,
		TelegramID: tgID,
	}, nil
}

func (r *repository) DeleteUserFromGroup(ctx context.Context, groupID int64, tgID int64) error {
	userID, err := r.GetUserIDByTelegramID(ctx, tgID)
	if err != nil {
		return fmt.Errorf("failed to get user id by telegram id: %w", err)
	}

	tx, err := r.storage.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	var rollback bool

	defer func() {
		if rollback {
			if err := tx.Rollback(ctx); err != nil {
				logger.GetLogger().Err(err).Msg("failed to rollback transaction")
			}
		}
	}()

	tag, err := tx.Exec(ctx, deleteUserFromGroupQuery, groupID, userID)
	if err != nil {
		rollback = true

		return fmt.Errorf("failed to execute delete user from group query: %w", err)
	}

	if tag.RowsAffected() != 1 {
		rollback = true

		return ErrUnexpectedResult
	}

	if err := tx.Commit(ctx); err != nil {
		rollback = true

		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *repository) getUserIDsWithAccessTypeInGroup(ctx context.Context, groupID int64, accessType model.RightsPolicy) ([]int64, error) {
	result, err := r.storage.Query(ctx, getUserIDsWithAccessTypeInGroupQuery, groupID, accessType)
	if err != nil {
		return nil, fmt.Errorf("failed to get user ids with access type in group query: %w", err)
	}

	ids := make([]int64, 0)

	for result.Next() {
		var id int64
		if err := result.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan user id: %w", err)
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (r *repository) GetUserIDsInGroupWithReadOnlyAccess(ctx context.Context, groupID int64) ([]int64, error) {
	result, err := r.getUserIDsWithAccessTypeInGroup(ctx, groupID, model.ReadOnlyRightPolicy)
	if err != nil {
		return nil, fmt.Errorf("failed to get user ids with access type in group: %w", err)
	}

	return result, nil
}

func (r *repository) GetUserIDsInGroupWithReadWriteAccess(ctx context.Context, groupID int64) ([]int64, error) {
	result, err := r.getUserIDsWithAccessTypeInGroup(ctx, groupID, model.ReadWriteRightPolicy)
	if err != nil {
		return nil, fmt.Errorf("failed to get user ids with access type in group: %w", err)
	}

	return result, nil
}

func (r *repository) ChangeMessage(ctx context.Context, chat_id int64, user_id int64, message_id int64) error {
	tag, err := r.storage.Exec(ctx, upsertMessage, chat_id, user_id, message_id)
	if err != nil {
		return fmt.Errorf("failed to execute upsert message query: %w", err)
	}

	if tag.RowsAffected() != 1 {
		return ErrUnexpectedResult
	}

	return nil
}
