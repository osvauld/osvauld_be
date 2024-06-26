package db

import (
	"context"
	"database/sql"
	dto "osvauld/dtos"

	"github.com/google/uuid"
)

type AddCredentialTransactionParams struct {
	Name                 string
	Description          sql.NullString
	FolderID             uuid.UUID
	CredentialType       string
	CreatedBy            uuid.UUID
	Fields               []dto.Fields
	CredentialAccessArgs []AddCredentialAccessParams
	Domain               string
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
			CreatedBy:      uuid.NullUUID{UUID: args.CreatedBy, Valid: true},
			CredentialType: args.CredentialType,
			Domain:         sql.NullString{String: args.Domain, Valid: true},
		}
		credentialID, err = q.CreateCredential(ctx, CreateCredentialParams)
		if err != nil {
			return err
		}

		// Create field records
		for _, field := range args.Fields {

			fieldID, err := q.AddFieldData(ctx, AddFieldDataParams{
				FieldName:    field.FieldName,
				FieldType:    field.FieldType,
				CredentialID: credentialID,
				CreatedBy:    uuid.NullUUID{UUID: args.CreatedBy, Valid: true},
			})
			if err != nil {
				return err
			}

			for _, userField := range field.FieldValues {
				_, err := q.AddFieldValue(ctx, AddFieldValueParams{
					FieldID:    fieldID,
					FieldValue: userField.FieldValue,
					UserID:     userField.UserID,
				})
				if err != nil {
					return err
				}
			}
		}

		for _, accessRow := range args.CredentialAccessArgs {

			accessRow.CredentialID = credentialID
			_, err := q.AddCredentialAccess(ctx, accessRow)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return credentialID, err
}

func (store *SQLStore) EditCredentialTransaction(ctx context.Context, args dto.EditCredentialRequest, editedBy uuid.UUID) error {

	err := store.execTx(ctx, func(q *Queries) error {

		var err error
		// Update credential record
		editCredentialDetailsParams := EditCredentialDetailsParams{
			ID:             args.CredentialID,
			Name:           args.Name,
			Description:    sql.NullString{String: args.Description, Valid: true},
			CredentialType: args.CredentialType,
			UpdatedBy:      uuid.NullUUID{UUID: editedBy, Valid: true},
			Domain:         sql.NullString{String: args.Domain, Valid: true},
		}
		err = q.EditCredentialDetails(ctx, editCredentialDetailsParams)
		if err != nil {
			return err
		}

		// Edit User Fields
		for _, field := range args.EditedUserFields {

			// Edit field data
			err = q.EditFieldData(ctx, EditFieldDataParams{
				ID:        field.FieldID,
				FieldName: field.FieldName,
				FieldType: field.FieldType,
				UpdatedBy: uuid.NullUUID{UUID: editedBy, Valid: true},
			})
			if err != nil {
				return err
			}

			for _, userValue := range field.FieldValues {
				err = q.EditFieldValue(ctx, EditFieldValueParams{
					FieldValue: userValue.FieldValue,
					UserID:     userValue.UserID,
					FieldID:    field.FieldID,
				})
				if err != nil {
					return err
				}
			}
		}

		// Edit Env Fields
		for _, field := range args.EditedEnvFields {
			err = q.EditEnvFieldValue(ctx, EditEnvFieldValueParams{
				FieldValue: field.FieldValue,
				ID:         field.EnvFieldID,
			})
			if err != nil {
				return err
			}
		}

		// Add new fields
		for _, field := range args.NewFields {

			// Add field data
			fieldID, err := q.AddFieldData(ctx, AddFieldDataParams{
				FieldName:    field.FieldName,
				FieldType:    field.FieldType,
				CredentialID: args.CredentialID,
				CreatedBy:    uuid.NullUUID{UUID: editedBy, Valid: true},
			})
			if err != nil {
				return err
			}

			for _, userField := range field.FieldValues {
				fieldValueID, err := q.AddFieldValue(ctx, AddFieldValueParams{
					FieldID:    fieldID,
					FieldValue: userField.FieldValue,
					UserID:     userField.UserID,
				})
				if err != nil {
					return err
				}

				for _, envField := range userField.EnvFieldValues {
					_, err = q.CreateEnvFields(ctx, CreateEnvFieldsParams{
						CredentialID:       args.CredentialID,
						FieldValue:         envField.FieldValue,
						FieldName:          field.FieldName,
						ParentFieldValueID: fieldValueID,
						EnvID:              envField.EnvID,
					})
					if err != nil {
						return err
					}
				}
			}

		}

		err = q.DeleteFields(ctx, args.DeletedFields)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
