package db

import (
	"context"
	"database/sql"
	dto "osvauld/dtos"

	"github.com/google/uuid"
)

type AddCredentialTransactionParams struct {
	Name                     string
	Description              sql.NullString
	FolderID                 uuid.UUID
	CredentialType           string
	CreatedBy                uuid.UUID
	UserFieldsWithAccessType []dto.UserFieldsWithAccessType
}

func (store *SQLStore) AddCredentialTransaction(ctx context.Context, args AddCredentialTransactionParams) (uuid.UUID, error) {

	var credentialID uuid.UUID

	err := store.execTx(ctx, func(q *Queries) error {

		var err error
		// Create credential record
		CreateCredentialParams := CreateCredentialParams{
			Name:           args.Name,
			Description:    args.Description,
			FolderID:       args.FolderID,
			CreatedBy:      args.CreatedBy,
			CredentialType: args.CredentialType,
		}
		credentialID, err = q.CreateCredential(ctx, CreateCredentialParams)
		if err != nil {
			return err
		}

		// Create field records
		for _, userFields := range args.UserFieldsWithAccessType {
			for _, field := range userFields.Fields {
				_, err = q.AddFieldData(ctx, AddFieldDataParams{
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

			accessListParams := AddCredentialAccessParams{
				CredentialID: credentialID,
				UserID:       userFields.UserID,
				AccessType:   userFields.AccessType,
			}
			q.AddCredentialAccess(ctx, accessListParams)
		}

		return nil
	})

	return credentialID, err
}
