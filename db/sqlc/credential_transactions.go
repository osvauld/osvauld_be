package db

import (
	"context"
	"database/sql"
	dto "osvauld/dtos"

	"github.com/google/uuid"
)

func (store *SQLStore) AddCredentialTransaction(ctx context.Context, args dto.AddCredentialDto) (uuid.UUID, error) {

	var credentialID uuid.UUID

	err := store.execTx(ctx, func(q *Queries) error {

		var err error
		// Create credential record
		CreateCredentialParams := CreateCredentialParams{
			Name:        args.Name,
			Description: sql.NullString{String: args.Description, Valid: true},
			FolderID:    args.FolderID,
			CreatedBy:   args.CreatedBy,
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
			})
			if err != nil {
				return err
			}
		}

		// Create encrypted data records
		for _, field := range args.EncryptedFields {
			_, err = q.CreateEncryptedData(ctx, CreateEncryptedDataParams{
				FieldName:    field.FieldName,
				FieldValue:   field.FieldValue,
				CredentialID: credentialID,
				UserID:       args.CreatedBy,
			})
			if err != nil {
				return err
			}
		}

		// Add rows in access list
		accessListParams := AddToAccessListParams{
			CredentialID: credentialID,
			UserID:       args.CreatedBy,
			AccessType:   "owner",
		}
		q.AddToAccessList(ctx, accessListParams)

		return nil
	})

	return credentialID, err
}
