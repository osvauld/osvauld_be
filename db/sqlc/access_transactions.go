package db

import (
	"context"

	"github.com/google/uuid"
)

func (store *SQLStore) RemoveFolderAccessForUsersTransactions(ctx context.Context, args RemoveFolderAccessForUsersParams) error {

	return store.execTx(ctx, func(q *Queries) error {

		// Remove folder access rows from folder_access table
		err := q.RemoveFolderAccessForUsers(ctx, args)
		if err != nil {
			return err
		}

		// Remove credential access rows from credential_access table for the folder
		err = q.RemoveCredentialAccessForUsersWithFolderID(ctx, RemoveCredentialAccessForUsersWithFolderIDParams{
			UserIds:  args.UserIds,
			FolderID: uuid.NullUUID{UUID: args.FolderID, Valid: true},
		})
		if err != nil {
			return err
		}

		return nil
	})
}

func (store *SQLStore) RemoveFolderAccessForGroupsTransactions(ctx context.Context, args RemoveFolderAccessForGroupsParams) error {
	return store.execTx(ctx, func(q *Queries) error {

		// Remove folder access rows from folder_access table
		err := q.RemoveFolderAccessForGroups(ctx, args)
		if err != nil {
			return err
		}

		// Remove credential access rows from credential_access table for the folder
		err = q.RemoveCredentialAccessForGroupsWithFolderID(ctx, RemoveCredentialAccessForGroupsWithFolderIDParams{
			GroupIds: args.GroupIds,
			FolderID: uuid.NullUUID{UUID: args.FolderID, Valid: true},
		})
		if err != nil {
			return err
		}

		return nil

	})

}

func (store *SQLStore) EditFolderAccessForUserTransaction(ctx context.Context, args EditFolderAccessForUserParams) error {
	return store.execTx(ctx, func(q *Queries) error {

		// Remove folder access rows from folder_access table
		err := q.EditFolderAccessForUser(ctx, EditFolderAccessForUserParams{
			AccessType: args.AccessType,
			UserID:     args.UserID,
			FolderID:   args.FolderID,
		})
		if err != nil {
			return err
		}

		// Add folder access rows to folder_access table
		err = q.EditCredentialAccessForUserWithFolderID(ctx, EditCredentialAccessForUserWithFolderIDParams{
			AccessType: args.AccessType,
			UserID:     args.UserID,
			FolderID:   uuid.NullUUID{UUID: args.FolderID, Valid: true},
		})
		if err != nil {
			return err
		}

		return nil
	})
}

// Edit folder access for group transaction
func (store *SQLStore) EditFolderAccessForGroupTransaction(ctx context.Context, args EditFolderAccessForGroupParams) error {
	return store.execTx(ctx, func(q *Queries) error {

		// Remove folder access rows from folder_access table
		err := q.EditFolderAccessForGroup(ctx, EditFolderAccessForGroupParams{
			AccessType: args.AccessType,
			GroupID:    args.GroupID,
			FolderID:   args.FolderID,
		})
		if err != nil {
			return err
		}

		// Add folder access rows to folder_access table
		err = q.EditCredentialAccessForGroupWithFolderID(ctx, EditCredentialAccessForGroupWithFolderIDParams{
			AccessType: args.AccessType,
			GroupID:    args.GroupID,
			FolderID:   uuid.NullUUID{UUID: args.FolderID, Valid: true},
		})
		if err != nil {
			return err
		}

		return nil
	})
}
