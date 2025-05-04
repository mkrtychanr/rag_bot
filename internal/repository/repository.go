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

type Repository struct {
	storage storage
}

func NewRepository(s storage) *Repository {
	return &Repository{
		storage: s,
	}
}

func (r *Repository) RegistrateUser(ctx context.Context, name string, tgID int64) (int64, error) {
	rows, err := r.storage.Query(ctx, registrateUserQuery, name, tgID)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return 0, ErrNoRows
	}

	var userID int64
	if err := rows.Scan(&userID); err != nil {
		return 0, fmt.Errorf("failed to scan userID: %w", err)
	}

	return userID, nil
}

func (r *Repository) GetUserIDByTelegramID(ctx context.Context, tgID int64) (int64, error) {
	rows, err := r.storage.Query(ctx, getUserIDByTelegramIDQuery, tgID)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return 0, ErrNoRows
	}

	var userID int64
	if err := rows.Scan(&userID); err != nil {
		return 0, fmt.Errorf("failed to scan userID: %w", err)
	}

	return userID, nil
}

func (r *Repository) CreateGroup(ctx context.Context, tgID int64, name string) error {
	userID, err := r.GetUserIDByTelegramID(ctx, tgID)
	if err != nil {
		return fmt.Errorf("failed to get user id by telegram id: %w", err)
	}

	rows, err := r.storage.Query(ctx, createGroupQuery, name, userID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return ErrNoRows
	}

	var groupID int64
	if err := rows.Scan(&groupID); err != nil {
		return fmt.Errorf("failed to scan groupID: %w", err)
	}

	if groupID == 0 {
		return ErrUnexpectedResult
	}

	return nil
}

func (r *Repository) AddUserIntoGroup(ctx context.Context, groupID int64, tgID int64) error {
	userID, err := r.GetUserIDByTelegramID(ctx, tgID)
	if err != nil {
		return fmt.Errorf("failed to get user id by telegram id: %w", err)
	}

	rows, err := r.storage.Query(ctx, addUserIntoGroupQuery, groupID, userID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return ErrNoRows
	}

	var groupUserID int64
	if err := rows.Scan(&groupUserID); err != nil {
		return fmt.Errorf("failed to scan groupUserID: %w", err)
	}

	if groupUserID == 0 {
		return ErrUnexpectedResult
	}

	return nil
}

func (r *Repository) SetReadOnlyRightsForUserInGroup(ctx context.Context, groupID int64, tgID int64) error {
	if err := r.updateUserRightsPolicy(ctx, groupID, tgID, model.ReadOnlyRightPolicy); err != nil {
		return fmt.Errorf("failed to update rights policy: %w", err)
	}

	return nil
}

func (r *Repository) SetReadWriteRightsForUserInGroup(ctx context.Context, groupID int64, tgID int64) error {
	if err := r.updateUserRightsPolicy(ctx, groupID, tgID, model.ReadWriteRightPolicy); err != nil {
		return fmt.Errorf("failed to update rights policy: %w", err)
	}

	return nil
}

func (r *Repository) updateUserRightsPolicy(ctx context.Context, groupID int64, userID int64, rightsPolicy model.RightsPolicy) error {
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

func (r *Repository) AddPaper(ctx context.Context, name string, tgID int64) (int64, error) {
	userID, err := r.GetUserIDByTelegramID(ctx, tgID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user id by telegram id: %w", err)
	}

	rows, err := r.storage.Query(ctx, addPaperQuery, name, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return 0, ErrNoRows
	}

	var paperID int64
	if err := rows.Scan(&paperID); err != nil {
		return 0, fmt.Errorf("failed to scan paperID: %w", err)
	}

	if paperID == 0 {
		return 0, ErrUnexpectedResult
	}

	return paperID, nil
}

func (r *Repository) AddPaperIntoGroup(ctx context.Context, paperID int64, groupID int64) error {
	rows, err := r.storage.Query(ctx, addPaperIntoGroupQuery, paperID, groupID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return ErrNoRows
	}

	var paperGroupID int64
	if err := rows.Scan(&paperGroupID); err != nil {
		return fmt.Errorf("failed to scan paperGroupID: %w", err)
	}

	if paperGroupID == 0 {
		return ErrUnexpectedResult
	}

	return nil
}

func (r *Repository) DeletePapers(ctx context.Context, papersIDs []int64) error {
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

func (r *Repository) DeleteGroups(ctx context.Context, groupIDs []int64) error {
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

func (r *Repository) GetUserGroups(ctx context.Context, tgID int64) ([]int64, error) {
	userID, err := r.GetUserIDByTelegramID(ctx, tgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user id by telegram id: %w", err)
	}

	rows, err := r.storage.Query(ctx, getUserGroupsQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getUserGroupsQuery: %w", err)
	}

	defer rows.Close()

	groups := make([]int64, 0)

	for rows.Next() {
		var group int64
		if err := rows.Scan(&group); err != nil {
			return nil, fmt.Errorf("failed to scan groups: %w", err)
		}

		groups = append(groups, group)
	}

	return groups, nil
}

func (r *Repository) GetUserRWGroupIDs(ctx context.Context, tgID int64) ([]int64, error) {
	userID, err := r.GetUserIDByTelegramID(ctx, tgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user id by telegram id: %w", err)
	}

	rows, err := r.storage.Query(ctx, getUserGroupIDsWithAccessTypeQuery, userID, model.ReadWriteRightPolicy)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get user group ids with access type query: %w", err)
	}

	defer rows.Close()

	groupIDs := make([]int64, 0)

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan group id: %w", err)
		}

		groupIDs = append(groupIDs, id)
	}

	return groupIDs, nil
}

func (r *Repository) FetchGroupsInfo(ctx context.Context, groupIDs []int64) ([]model.Group, error) {
	if len(groupIDs) == 0 {
		return nil, nil
	}
	rows, err := r.storage.Query(ctx, getGroupsInfoQuery, groupIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get groups info query: %w", err)
	}

	defer rows.Close()

	groups := make([]model.Group, 0, len(groupIDs))

	for rows.Next() {
		var (
			id      int64
			name    string
			adminID int64
		)

		if err := rows.Scan(&id, &name, &adminID); err != nil {
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

func (r *Repository) GetUserPapers(ctx context.Context, tgID int64) ([]model.Paper, error) {
	userID, err := r.GetUserIDByTelegramID(ctx, tgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user id by telegram id: %w", err)
	}

	rows, err := r.storage.Query(ctx, getUserPapersQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get user papers query: %w", err)
	}

	defer rows.Close()

	papers := make([]model.Paper, 0)

	for rows.Next() {
		var (
			id   int64
			name string
		)

		if err := rows.Scan(&id, &name); err != nil {
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

func (r *Repository) DeletePaperFromGroup(ctx context.Context, paperID int64, groupID int64) error {
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

func (r *Repository) GetUserGroupsOwnership(ctx context.Context, tgID int64) ([]model.Group, error) {
	userID, err := r.GetUserIDByTelegramID(ctx, tgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user id by telegram id: %w", err)
	}

	rows, err := r.storage.Query(ctx, getUserGroupsOwnershipQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get user groups ownership query: %w", err)
	}

	defer rows.Close()

	groups := make([]model.Group, 0)

	for rows.Next() {
		var (
			id   int64
			name string
		)

		if err := rows.Scan(&id, &name); err != nil {
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

func (r *Repository) GetGroupPapers(ctx context.Context, groupID int64) ([]int64, error) {
	rows, err := r.storage.Query(ctx, getGroupPapersQuery, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get group papers query: %w", err)
	}

	defer rows.Close()

	papers := make([]int64, 0)

	for rows.Next() {
		var id int64

		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan paper id: %w", err)
		}

		papers = append(papers, id)
	}

	return papers, nil
}

func (r *Repository) FetchPapersInfo(ctx context.Context, paperIDs []int64) ([]model.Paper, error) {
	rows, err := r.storage.Query(ctx, getPapersInfoQuery, paperIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get papers info query: %w", err)
	}

	defer rows.Close()

	papers := make([]model.Paper, 0, len(paperIDs))

	for rows.Next() {
		var (
			id       int64
			name     string
			authorID int64
		)

		if err := rows.Scan(&id, &name, &authorID); err != nil {
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

func (r *Repository) GetGroupUsers(ctx context.Context, groupID int64) ([]model.UserGroup, error) {
	rows, err := r.storage.Query(ctx, getGroupUsersQuery, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get group users query: %w", err)
	}

	defer rows.Close()

	users := make(map[int64]model.UserGroup)

	for rows.Next() {
		var (
			id         int64
			accessType model.RightsPolicy
		)

		if err := rows.Scan(&id, &accessType); err != nil {
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

func (r *Repository) FetchUsersInfo(ctx context.Context, userIDs []int64) ([]model.User, error) {
	rows, err := r.storage.Query(ctx, getUsersInfoQuery, userIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get users info query: %w", err)
	}

	defer rows.Close()

	users := make([]model.User, 0)

	for rows.Next() {
		var (
			id   int64
			name string
			tgID int64
		)

		if err := rows.Scan(&id, &name, &tgID); err != nil {
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

func (r *Repository) ChangeGroupName(ctx context.Context, groupID int64, name string) error {
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

func (r *Repository) GetUserByShortname(ctx context.Context, shortname string) (model.User, error) {
	rows, err := r.storage.Query(ctx, getUserByShortnameQuery, shortname)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to execute get user by shortname query: %w", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return model.User{}, ErrNoRows
	}

	var (
		id   int64
		name string
		tgID int64
	)

	if err := rows.Scan(&id, &name, &tgID); err != nil {
		return model.User{}, fmt.Errorf("failed to scan user info: %w", err)
	}

	return model.User{
		ID:         id,
		Shortname:  name,
		TelegramID: tgID,
	}, nil
}

func (r *Repository) DeleteUserFromGroup(ctx context.Context, groupID int64, tgID int64) error {
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

func (r *Repository) getUserIDsWithAccessTypeInGroup(ctx context.Context, groupID int64, accessType model.RightsPolicy) ([]int64, error) {
	rows, err := r.storage.Query(ctx, getUserIDsWithAccessTypeInGroupQuery, groupID, accessType)
	if err != nil {
		return nil, fmt.Errorf("failed to get user ids with access type in group query: %w", err)
	}

	defer rows.Close()

	ids := make([]int64, 0)

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan user id: %w", err)
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (r *Repository) GetUserIDsInGroupWithReadOnlyAccess(ctx context.Context, groupID int64) ([]int64, error) {
	result, err := r.getUserIDsWithAccessTypeInGroup(ctx, groupID, model.ReadOnlyRightPolicy)
	if err != nil {
		return nil, fmt.Errorf("failed to get user ids with access type in group: %w", err)
	}

	return result, nil
}

func (r *Repository) GetUserIDsInGroupWithReadWriteAccess(ctx context.Context, groupID int64) ([]int64, error) {
	result, err := r.getUserIDsWithAccessTypeInGroup(ctx, groupID, model.ReadWriteRightPolicy)
	if err != nil {
		return nil, fmt.Errorf("failed to get user ids with access type in group: %w", err)
	}

	return result, nil
}

func (r *Repository) ChangeMessage(ctx context.Context, chatID int64, userID int64, messageID int64) error {
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

	tag, err := tx.Exec(ctx, upsertMessageQuery, chatID, userID, messageID)
	if err != nil {
		rollback = true

		return fmt.Errorf("failed to execute upsert message query: %w", err)
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

func (r *Repository) GetMessageID(ctx context.Context, chatID int64) (int64, error) {
	rows, err := r.storage.Query(ctx, getMessageIDQuery, chatID)
	if err != nil {
		return 0, fmt.Errorf("failed to execute get message id query: %w", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return 0, ErrNoRows
	}

	var id int64
	if err := rows.Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to scan message id: %w", err)
	}

	return id, nil
}

func (r *Repository) DeleteMessage(ctx context.Context, chatID int64) error {
	tag, err := r.storage.Exec(ctx, deleteMessageQuery, chatID)
	if err != nil {
		return fmt.Errorf("failed to execute delete message query: %w", err)
	}

	if tag.RowsAffected() != 1 {
		return ErrUnexpectedResult
	}

	return nil
}

func (r *Repository) GetPaperGroupIDs(ctx context.Context, paperID int64) ([]int64, error) {
	rows, err := r.storage.Query(ctx, getPaperGroupIDsQuery, paperID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get paper group ids query: %w", err)
	}

	defer rows.Close()

	ids := make([]int64, 0)

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan group id: %w", err)
		}

		ids = append(ids, id)
	}

	return ids, nil
}
