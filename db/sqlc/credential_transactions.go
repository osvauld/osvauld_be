package db

import (
	"context"
	"database/sql"

	dto "osvauld/dtos"

	"github.com/google/uuid"
)

type AddCredentialTransactionParams struct {
	Name              string                  `json:"name"`
	Description       string                  `json:"description"`
	FolderID          uuid.UUID               `json:"folderId"`
	UnencryptedFields []dto.Field             `json:"unencryptedFields"`
	UserAccessDetails []dto.UserEncryptedData `json:"userAccessDetails"`
	CreatedBy         uuid.UUID               `json:"createdBy"`
}

func (store *SQLStore) AddCredentialTransaction(ctx context.Context, args AddCredentialTransactionParams) (uuid.UUID, error) {

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
		for _, userDetails := range args.UserAccessDetails {
			for _, field := range userDetails.EncryptedFields {
				_, err = q.CreateEncryptedData(ctx, CreateEncryptedDataParams{
					FieldName:    field.FieldName,
					FieldValue:   field.FieldValue,
					CredentialID: credentialID,
					UserID:       userDetails.UserID,
				})
				if err != nil {
					return err
				}
			}
		}

		// Add rows in access list
		for _, user := range args.UserAccessDetails {
			accessListParams := AddToAccessListParams{
				CredentialID: credentialID,
				UserID:       user.UserID,
				AccessType:   user.AccessType,
			}
			q.AddToAccessList(ctx, accessListParams)
		}

		return nil
	})

	return credentialID, err
}

type ShareCredentialParams struct {
	CredentialID      uuid.UUID               `json:"credentialId"`
	UserEncryptedData []dto.UserEncryptedData `json:"encryptedFields"`
}

func (store *SQLStore) ShareCredentialTransaction(ctx context.Context, args ShareCredentialParams) error {

	err := store.execTx(ctx, func(q *Queries) error {

		// Create encrypted data records
		for _, userDetails := range args.UserEncryptedData {
			for _, field := range userDetails.EncryptedFields {
				_, err := q.CreateEncryptedData(ctx, CreateEncryptedDataParams{
					FieldName:    field.FieldName,
					FieldValue:   field.FieldValue,
					CredentialID: args.CredentialID,
					UserID:       userDetails.UserID,
				})
				if err != nil {
					return err
				}
			}
		}

		// Add rows in access list
		for _, user := range args.UserEncryptedData {
			accessListParams := AddToAccessListParams{
				CredentialID: args.CredentialID,
				UserID:       user.UserID,
				AccessType:   user.AccessType,
				GroupID:      user.GroupID,
			}
			q.AddToAccessList(ctx, accessListParams)
		}

		return nil
	})

	return err
}
