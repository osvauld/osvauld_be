package db

import (
	"context"
)

type ShareCredentialTransactionParams struct {
	FieldArgs            []AddFieldValueParams       `json:"fieldArgs"`
	CredentialAccessArgs []AddCredentialAccessParams `json:"credentialAccessArgs"`
	FolderAccessArgs     []AddFolderAccessParams     `json:"folderAccessArgs"`
}

func (store *SQLStore) ShareCredentialsTransaction(ctx context.Context, args ShareCredentialTransactionParams) error {

	err := store.execTx(ctx, func(q *Queries) error {

		for _, fieldRecord := range args.FieldArgs {

			err := q.AddFieldValue(ctx, fieldRecord)
			if err != nil {
				return err
			}
		}

		for _, credentialAccessRecord := range args.CredentialAccessArgs {

			_, err := q.AddCredentialAccess(ctx, credentialAccessRecord)
			if err != nil {
				return err
			}

		}

		for _, folderAccessRecord := range args.FolderAccessArgs {

			err := q.AddFolderAccess(ctx, folderAccessRecord)
			if err != nil {
				return err
			}

		}

		return nil
	})

	return err
}
