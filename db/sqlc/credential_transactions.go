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
				_, err = q.AddField(ctx, AddFieldParams{
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

type EditCredentialTransactionParams struct {
	CredentialID   uuid.UUID
	Name           string
	Description    sql.NullString
	CredentialType string
	UpdatedBy      uuid.UUID
	EditFields     []dto.UserFields
	AddFields      []dto.UserFieldsWithAccessType
}

func (store *SQLStore) EditCredentialTransaction(ctx context.Context, args EditCredentialTransactionParams) error {

	err := store.execTx(ctx, func(q *Queries) error {

		var err error
		// Update credential record
		editCredentialDetailsParams := EditCredentialDetailsParams{
			ID:             args.CredentialID,
			Name:           args.Name,
			Description:    args.Description,
			CredentialType: args.CredentialType,
		}
		err = q.EditCredentialDetails(ctx, editCredentialDetailsParams)
		if err != nil {
			return err
		}

		// Edit field records
		for _, userFields := range args.EditFields {
			for _, field := range userFields.Fields {
				err = q.EditField(ctx, EditFieldParams{
					ID:         field.ID,
					FieldName:  field.FieldName,
					FieldValue: field.FieldValue,
					FieldType:  field.FieldType,
				})
				if err != nil {
					return err
				}
			}
		}

		// Create field records
		for _, userFields := range args.AddFields {
			for _, field := range userFields.Fields {
				_, err = q.AddField(ctx, AddFieldParams{
					FieldName:    field.FieldName,
					FieldValue:   field.FieldValue,
					CredentialID: args.CredentialID,
					UserID:       userFields.UserID,
					FieldType:    field.FieldType,
				})
				if err != nil {
					return err
				}
			}

			accessListParams := AddCredentialAccessParams{
				CredentialID: args.CredentialID,
				UserID:       userFields.UserID,
				AccessType:   userFields.AccessType,
			}
			q.AddCredentialAccess(ctx, accessListParams)
		}

		return nil
	})

	return err
}
