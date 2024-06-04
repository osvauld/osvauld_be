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
				err := q.AddFieldValue(ctx, AddFieldValueParams{
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

type EditCredentialTransactionParams struct {
	CredentialID   uuid.UUID
	Name           string
	Description    sql.NullString
	CredentialType string
	EditedFields   []dto.Fields
	AddFields      []dto.Fields
	EditedBy       uuid.UUID
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

		// // Create field records
		// for _, field := range args.NewFields {

		// 	// Add field data
		// 	fieldID, err := q.AddFieldData(ctx, AddFieldDataParams{
		// 		FieldName:    field.FieldName,
		// 		FieldType:    field.FieldType,
		// 		CredentialID: args.CredentialID,
		// 		CreatedBy:    uuid.NullUUID{UUID: args.EditedBy, Valid: true},
		// 	})
		// 	if err != nil {
		// 		return err
		// 	}

		// 	// Add field values
		// 	for _, userField := range field.FieldValues {

		// 		isCliUser, err := q.CheckCliUser(ctx, userField.UserID)
		// 		if err != nil {
		// 			return err
		// 		}

		// 		if isCliUser {
		// 			envs, err := q.GetUserEnvsForCredential(ctx, GetUserEnvsForCredentialParams{
		// 				CredentialID: args.CredentialID,
		// 				CliUser:      userField.UserID,
		// 			})
		// 			if err != nil {
		// 				return err
		// 			}

		// 			for _, envID := range envs {
		// 				_, err := q.CreateEnvFields(ctx, CreateEnvFieldsParams{
		// 					CredentialID:  args.CredentialID,
		// 					FieldValue:    userField.FieldValue,
		// 					FieldName:     field.FieldName,
		// 					ParentFieldID: fieldID,
		// 					EnvID:         envID,
		// 				})
		// 				if err != nil {
		// 					return err
		// 				}
		// 			}
		// 		} else {
		// 			err := q.AddFieldValue(ctx, AddFieldValueParams{
		// 				FieldID:    fieldID,
		// 				FieldValue: userField.FieldValue,
		// 				UserID:     userField.UserID,
		// 			})
		// 			if err != nil {
		// 				return err
		// 			}
		// 		}
		// 	}

		// }
		

		return nil
	})

	return err
}
