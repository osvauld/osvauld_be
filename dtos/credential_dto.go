package dto

import (
	db "osvauld/db/sqlc"

	"github.com/google/uuid"
)

type AddCredentailRequest struct {
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	FolderID          uuid.UUID      `json:"folderId"`
	EncryptedFields   []FieldRequest `json:"encryptedFields"`
	UnencryptedFields []FieldRequest `json:"unencryptedFields"`
}

type FieldRequest struct {
	FieldName  string `json:"fieldName"`
	FieldValue string `json:"fieldValue"`
}

type User struct {
	UserID     uuid.UUID      `json:"userId"`
	Fields     []FieldRequest `json:"fields"`
	AccessType string         `json:"accessType"`
}

type Credential struct {
	CredentialID uuid.UUID `json:"credentialId"`
	Users        []User    `json:"users"`
}

type ShareCredentialPayload struct {
	CredentialList []Credential `json:"credentialList"`
}

type CredentialDetails struct {
	Credential      db.GetCredentialDetailsRow           `json:"credential"`
	EncryptedData   []db.GetUserEncryptedDataRow         `json:"encryptedData"`
	UnencryptedData []db.GetCredentialUnencryptedDataRow `json:"unencryptedData"`
}
