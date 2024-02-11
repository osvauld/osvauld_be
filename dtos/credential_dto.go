package dto

import (
	"time"

	"github.com/google/uuid"
)

type AddCredentialRequest struct {
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	FolderID       uuid.UUID    `json:"folderId"`
	CredentialType string       `json:"credentialType"`
	UserFields     []UserFields `json:"userFields"`
}

type AddCredentialDto struct {
	Name                     string                     `json:"name"`
	Description              string                     `json:"description"`
	FolderID                 uuid.UUID                  `json:"folderId"`
	CredentialType           string                     `json:"Credentialtype"`
	UserFieldsWithAccessType []UserFieldsWithAccessType `json:"userFieldsWithAccessType"`
}

type UserEncryptedData struct {
	UserID          uuid.UUID     `json:"userId"`
	AccessType      string        `json:"accessType"`
	GroupID         uuid.NullUUID `json:"groupId"`
	EncryptedFields []Field       `json:"encryptedFields"`
}

type CredentialForUser struct {
	CredentialID   uuid.UUID `json:"credentialId"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	FolderID       uuid.UUID `json:"folderId"`
	CredentialType string    `json:"credentialType"`
	AccessType     string    `json:"accessType"`
	Fields         []Field   `json:"fields"`
	CreatedBy      uuid.UUID `json:"createdBy"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type CredentialDetails struct {
	CredentialID   uuid.UUID `json:"credentialId"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	FolderID       uuid.UUID `json:"folderId"`
	CredentialType string    `json:"credentialType"`
	AccessType     string    `json:"accessType"`
	Fields         []Field   `json:"encryptedFields"`
	UserID         uuid.UUID `json:"userId"`
	CreatedBy      uuid.UUID `json:"createdBy"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
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

type GetCredentialsFieldsByIdsRequest struct {
	CredentialIds []uuid.UUID `json:"credentialIds"`
}

type GetCredentialsByIDsRequest struct {
	CredentialIds []uuid.UUID `json:"credentialIds"`
}
