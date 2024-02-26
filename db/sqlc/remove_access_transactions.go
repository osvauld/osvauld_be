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
