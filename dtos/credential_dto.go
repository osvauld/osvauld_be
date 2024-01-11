package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserEncryptedData struct {
	UserID          uuid.UUID     `json:"userId"`
	AccessType      string        `json:"accessType"`
	GroupID         uuid.NullUUID `json:"groupId"`
	EncryptedFields []Field       `json:"encryptedFields"`
}

type CredentialDetails struct {
	CredentialID      uuid.UUID `json:"credentialId"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	FolderID          uuid.UUID `json:"folderId"`
	UnencryptedFields []Field   `json:"unencryptedFields"`
	EncryptedFields   []Field   `json:"encryptedFields"`
	UserID            uuid.UUID `json:"userId"`
	CreatedBy         uuid.UUID `json:"createdBy"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type Field struct {
	ID         uuid.UUID `json:"id"`
	FieldName  string    `json:"fieldName"`
	FieldValue string    `json:"fieldValue"`
}

type User struct {
	UserID     uuid.UUID `json:"userId"`
	Fields     []Field   `json:"fields"`
	AccessType string    `json:"accessType"`
}

type Credential struct {
	CredentialID uuid.UUID `json:"credentialId"`
	Users        []User    `json:"users"`
}

type GetEncryptedCredentialsByIdsRequest struct {
	CredentialIds []uuid.UUID `json:"credentialIds"`
}

type CredentialsForUser struct {
	CredentialID uuid.UUID `json:"credentialId"`
}
