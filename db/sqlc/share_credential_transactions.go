package db

import (
	"context"

	"github.com/google/uuid"
)

type FieldRecord struct {
	FieldName    string    `json:"fieldName"`
	FieldValue   string    `json:"fieldValue"`
	FieldType    string    `json:"fieldType"`
	CredentialID uuid.UUID `json:"credentialId"`
	UserID       uuid.UUID `json:"userId"`
}

type CredentialAccessRecord struct {
	CredentialID uuid.UUID     `json:"credentialId"`
	UserID       uuid.UUID     `json:"userId"`
	AccessType   string        `json:"accessType"`
	GroupID      uuid.NullUUID `json:"groupId"`
}

type FolderAccessRecord struct {
	FolderID   uuid.UUID     `json:"folderId"`
	UserID     uuid.UUID     `json:"userId"`
	AccessType string        `json:"accessType"`
	GroupID    uuid.NullUUID `json:"groupId"`
}

type ShareCredentialTransactionParams struct {
	FieldArgs            []AddFieldDataParams        `json:"fieldArgs"`
	CredentialAccessArgs []AddCredentialAccessParams `json:"credentialAccessArgs"`
	FolderAccessArgs     []AddFolderAccessParams     `json:"folderAccessArgs"`
}

func (store *SQLStore) ShareCredentialsTransaction(ctx context.Context, args ShareCredentialTransactionParams) error {

	err := store.execTx(ctx, func(q *Queries) error {

		for _, fieldRecord := range args.FieldArgs {

			_, err := q.AddFieldData(ctx, fieldRecord)
			if err != nil {
				return err
			}
		}

		for _, credentialAccessRecord := range args.CredentialAccessArgs {

			_, err := q.AddCredentialAccess(ctx, credentialAccessRecord)
			if err != nil {
				return err
			}

		}

		for _, folderAccessRecord := range args.FolderAccessArgs {

			err := q.AddFolderAccess(ctx, folderAccessRecord)
			if err != nil {
				return err
			}

		}

		return nil
	})

	return err
}
