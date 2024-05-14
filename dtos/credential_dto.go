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
	Domain         string       `json:"domain"`
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

type GetCredentialsFieldsByIdsRequest struct {
	CredentialIds []uuid.UUID `json:"credentialIds"`
}

type GetCredentialsByIDsRequest struct {
	CredentialIds []uuid.UUID `json:"credentialIds"`
}

type ShareCredentialPayload struct {
	CredentialID uuid.UUID    `json:"credentialId" binding:"required"`
	Fields       []ShareField `json:"fields" binding:"required"`
}

type ShareCredentialsForUserPayload struct {
	UserID         uuid.UUID                `json:"userId" binding:"required"`
	AccessType     string                   `json:"accessType"`
	CredentialData []ShareCredentialPayload `json:"credentials" binding:"required"`
}

type CredentialsForGroupsPayload struct {
	GroupID    uuid.UUID                        `json:"groupId" binding:"required"`
	AccessType string                           `json:"accessType" binding:"required"`
	UserData   []ShareCredentialsForUserPayload `json:"userData" binding:"required"`
}

type ShareCredentialsWithUsersRequest struct {
	UserData []ShareCredentialsForUserPayload `json:"userData" binding:"required"`
}

type ShareCredentialsWithGroupsRequest struct {
	GroupData []CredentialsForGroupsPayload `json:"groupData" binding:"required"`
}

type ShareFolderWithUsersRequest struct {
	FolderID uuid.UUID                        `json:"folderId" binding:"required"`
	UserData []ShareCredentialsForUserPayload `json:"userData" binding:"required"`
}

type ShareFolderWithGroupsRequest struct {
	FolderID  uuid.UUID                     `json:"folderId" binding:"required"`
	GroupData []CredentialsForGroupsPayload `json:"groupData" binding:"required"`
}

type EditCredentialRequest struct {
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	CredentialType string       `json:"credentialType"`
	UserFields     []UserFields `json:"userFields"`
}


type EnvField struct {
	FieldID    uuid.UUID `json:"fieldId"`
	FieldValue string    `json:"fieldValue"`
}

type CredentialsData struct {
	CredentialID uuid.UUID  `json:"credentialId"`
	Fields       []EnvField `json:"fields"`
}

type ShareCredentialsWithEnvironmentRequest struct {
	EnvId       uuid.UUID         `json:"envId"`
	Credentials []CredentialsData `json:"credentials"`
}

type CredentialEnvData struct {
	CredentialID  uuid.UUID
	EnvID         uuid.UUID
	ParentFieldId uuid.UUID
	FieldValue    string
	FieldName     string
	CliUser       uuid.UUID
}
type EditCredentialDetailsRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	CredentialType string `json:"credentialType"`

}
