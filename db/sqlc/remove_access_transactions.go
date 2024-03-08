package db

import (
	"context"

	"github.com/google/uuid"
)

func (store *SQLStore) RemoveCredentialAccessForUsersTransactions(ctx context.Context, args RemoveCredentialAccessForUsersParams) error {
	return store.execTx(ctx, func(q *Queries) error {

		// Remove access from access list
		err := q.RemoveCredentialAccessForUsers(ctx, args)
		if err != nil {
			return err
		}

		for _, userID := range args.UserIds {

			exists, err := q.CheckCredentialAccessEntryExists(ctx, CheckCredentialAccessEntryExistsParams{
				UserID:       userID,
				CredentialID: args.CredentialID,
			})
			if err != nil {
				return err
			}

			// delete the fields if the user has no other access to the credential
			if !exists {

				err = q.RemoveCredentialFieldsForUsers(ctx, RemoveCredentialFieldsForUsersParams{
					UserIds:      []uuid.UUID{userID},
					CredentialID: args.CredentialID,
				})
				if err != nil {
					return err
				}
			}
		}

		return nil

	})
}

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

		// TODO: field cleanup
		return nil
	})
}

func (store *SQLStore) RemoveCredentialAccessForGroupsTransactions(ctx context.Context, args RemoveCredentialAccessForGroupsParams) error {
	return store.execTx(ctx, func(q *Queries) error {

		// Remove access from access list
		err := q.RemoveCredentialAccessForGroups(ctx, args)
		if err != nil {
			return err
		}

		// TODO: field cleanup

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

		// TODO: field cleanup

		return nil

	})

}
