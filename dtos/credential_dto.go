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

type EditCredentialDetailsRequest struct {
	CredentialID   uuid.UUID `json:"credentialId"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	CredentialType string    `json:"credentialType"`
}

type EditCredentialRequest struct {
	CredentialID   uuid.UUID                  `json:"credentialId"`
	Name           string                     `json:"name"`
	Description    string                     `json:"description"`
	CredentialType string                     `json:"credentialType"`
	EditFields     []UserFields               `json:"editFields"`
	AddFields      []UserFieldsWithAccessType `json:"addFields"`
}
