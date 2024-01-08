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
	UserAccessDetails []UserEncryptedData `json:"userAccessDetails"`
}

type UserEncryptedData struct {
	UserID          uuid.UUID     `json:"userId"`
	AccessType      string        `json:"accessType"`
	GroupID         uuid.NullUUID `json:"groupId"`
	EncryptedFields []Field       `json:"encryptedFields"`
}

type UserEncryptedFields struct {
	UserID          uuid.UUID `json:"userId"`
	EncryptedFields []Field   `json:"encryptedFields"`
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

type ShareCredentialWithUsers struct {
	CredentialID      uuid.UUID           `json:"credentialId"`
	UserEncryptedData []UserEncryptedData `json:"userEncryptedData"`
}

type ShareCredentialWithGroups struct {
	CredentialID uuid.UUID                `json:"credentialId"`
	GroupData    []CredentialDataForGroup `json:"groupData"`
}

type ShareMultipleCredentialsWithMultipleUsersPayload struct {
	Credentials []ShareCredentialWithUsers `json:"credentials"`
}

type ShareMultipleCredentialsWithMultipleGroupsPayload struct {
	Credentials []ShareCredentialWithGroups `json:"credentials"`
}

type GetEncryptedCredentialsByIdsRequest struct {
	CredentialIds []uuid.UUID `json:"credentialIds"`
}

type CredentialDataForGroup struct {
	GroupID             uuid.UUID             `json:"groupId"`
	UserEncryptedFields []UserEncryptedFields `json:"userEncryptedFields"`
	AccessType          string                `json:"accessType"`
}
