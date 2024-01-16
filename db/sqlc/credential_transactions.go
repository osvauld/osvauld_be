package db

import (
	"context"
	"database/sql"
	dto "osvauld/dtos"

	"github.com/google/uuid"
)

func (store *SQLStore) AddCredentialTransaction(ctx context.Context, args dto.AddCredentialDto, caller uuid.UUID) (uuid.UUID, error) {

	var credentialID uuid.UUID

	err := store.execTx(ctx, func(q *Queries) error {

		var err error
		// Create credential record
		CreateCredentialParams := CreateCredentialParams{
			Name:        args.Name,
			Description: sql.NullString{String: args.Description, Valid: true},
			FolderID:    args.FolderID,
			CreatedBy:   caller,
		}
		credentialID, err = q.CreateCredential(ctx, CreateCredentialParams)
		if err != nil {
			return err
		}

		// Create unencrypted data records
		for _, field := range args.UnencryptedFields {
			_, err = q.CreateUnencryptedData(ctx, CreateUnencryptedDataParams{
				FieldName:    field.FieldName,
				FieldValue:   field.FieldValue,
				CredentialID: credentialID,
				Url:          sql.NullString{String: field.URL, Valid: true},
				IsUrl:        field.IsUrl,
			})
			if err != nil {
				return err
			}
		}

		// Create encrypted data records
		for _, userEncryptedFields := range args.UserEncryptedFieldsWithAccess {
			for _, field := range userEncryptedFields.EncryptedFields {
				_, err = q.CreateEncryptedData(ctx, CreateEncryptedDataParams{
					FieldName:    field.FieldName,
					FieldValue:   field.FieldValue,
					CredentialID: credentialID,
					UserID:       userEncryptedFields.UserID,
				})
				if err != nil {
					return err
				}
			}

			accessListParams := AddToAccessListParams{
				CredentialID: credentialID,
				UserID:       userEncryptedFields.UserID,
				AccessType:   userEncryptedFields.AccessType,
			}
			q.AddToAccessList(ctx, accessListParams)
		}

		// Add rows in access list

		return nil
	})

	return credentialID, err
}
