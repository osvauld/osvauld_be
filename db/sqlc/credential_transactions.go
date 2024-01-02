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

type ShareCredentialWithUserParams struct {
	CredentialID    uuid.UUID     `json:"credentialId"`
	UserID          uuid.UUID     `json:"userId"`
	EncryptedFields []dto.Field   `json:"encryptedFields"`
	AccessType      string        `json:"accessType"`
	GroupID         uuid.NullUUID `json:"groupId"`
}

func (store *SQLStore) ShareCredentialWithUserTransaction(ctx context.Context, args ShareCredentialWithUserParams) error {

	err := store.execTx(ctx, func(q *Queries) error {

		// Create encrypted data records
		for _, field := range args.EncryptedFields {
			_, err := q.CreateEncryptedData(ctx, CreateEncryptedDataParams{
				FieldName:    field.FieldName,
				FieldValue:   field.FieldValue,
				CredentialID: args.CredentialID,
				UserID:       args.UserID,
			})
			if err != nil {
				return err
			}
		}

		// Add row in access list
		accessListParams := AddToAccessListParams{
			CredentialID: args.CredentialID,
			UserID:       args.UserID,
			AccessType:   args.AccessType,
			GroupID:      args.GroupID,
		}
		_, err := q.AddToAccessList(ctx, accessListParams)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
