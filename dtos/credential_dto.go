package dto

import (
	"time"

	"github.com/google/uuid"
)

type AddCredentailRequest struct {
	Name              string              `json:"name"`
	Description       string              `json:"description"`
	FolderID          uuid.UUID           `json:"folderId"`
	UnencryptedFields []Field             `json:"unencryptedFields"`
	UserAccessDetails []UserAccessDetails `json:"userAccessDetails"`
}

type UserAccessDetails struct {
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

type ShareCredentialPayload struct {
	CredentialList []Credential `json:"credentialList"`
}

// type CredentialDetails struct {
// 	Credential      db.GetCredentialDetailsRow           `json:"credential"`
// 	EncryptedData   []EncryptedFields        `json:"encryptedData"`
// 	UnencryptedData []db.GetCredentialUnencryptedDataRow `json:"unencryptedData"`
// 	Users           []db.GetUsersByCredentialRow         `json:"users"`
// }

type GetEncryptedCredentialsByIdsRequest struct {
	CredentialIds []uuid.UUID `json:"credentialIds"`
}
