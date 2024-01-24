package db

import (
	"context"
	dto "osvauld/dtos"
)

func (store *SQLStore) ShareCredentialsTransaction(ctx context.Context, args dto.ShareCredentialTransactionParams) error {

	err := store.execTx(ctx, func(q *Queries) error {

		for _, credentialData := range args.CredentialArgs {

			// Create encrypted data records
			for _, field := range credentialData.Fields {
				_, err := q.CreateFieldData(ctx, CreateFieldDataParams{
					FieldName:    field.FieldName,
					FieldValue:   field.FieldValue,
					FieldType:    field.FieldType,
					CredentialID: credentialData.CredentialID,
					UserID:       credentialData.UserID,
				})
				if err != nil {
					return err
				}
			}

			// Add row in access list
			accessListParams := AddToAccessListParams{
				CredentialID: credentialData.CredentialID,
				UserID:       credentialData.UserID,
				AccessType:   credentialData.AccessType,
				GroupID:      credentialData.GroupID,
			}
			_, err := q.AddToAccessList(ctx, accessListParams)
			if err != nil {
				return err
			}

		}

		for _, folderAccessData := range args.FolderAccess {

			// Add row in access list
			accessListParams := AddFolderAccessWithGroupParams{
				FolderID:   folderAccessData.FolderID,
				UserID:     folderAccessData.UserID,
				AccessType: folderAccessData.AccessType,
				GroupID:    folderAccessData.GroupID,
			}

			err := q.AddFolderAccessWithGroup(ctx, accessListParams)
			if err != nil {
				return err
			}

		}

		return nil
	})

	return err
}
