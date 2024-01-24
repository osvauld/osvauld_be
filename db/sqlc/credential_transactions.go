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
			Name:           args.Name,
			Description:    sql.NullString{String: args.Description, Valid: true},
			FolderID:       args.FolderID,
			CreatedBy:      caller,
			CredentialType: args.CredentialType,
		}
		credentialID, err = q.CreateCredential(ctx, CreateCredentialParams)
		if err != nil {
			return err
		}

		// Create field records
		for _, userFields := range args.UserFieldsWithAccessType {
			for _, field := range userFields.Fields {
				_, err = q.CreateFieldData(ctx, CreateFieldDataParams{
					FieldName:    field.FieldName,
					FieldValue:   field.FieldValue,
					CredentialID: credentialID,
					UserID:       userFields.UserID,
					FieldType:    field.FieldType,
				})
				if err != nil {
					return err
				}
			}

			accessListParams := AddToAccessListParams{
				CredentialID: credentialID,
				UserID:       userFields.UserID,
				AccessType:   userFields.AccessType,
			}
			q.AddToAccessList(ctx, accessListParams)
		}

		return nil
	})

	return credentialID, err
}
