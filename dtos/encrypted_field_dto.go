package dto

import "github.com/google/uuid"

type UserEncryptedFieldsDto struct {
	UserID          uuid.UUID `json:"userId"`
	EncryptedFields []Field   `json:"encryptedFields"`
}

type CredentialEncryptedFieldsForUserDto struct {
	CredentialID    uuid.UUID `json:"credentialId"`
	UserID          uuid.UUID `json:"userId"`
	EncryptedFields []Field   `json:"encryptedFields"`
	AccessType      string    `json:"accessType"`
}

type CredentialEncryptedFieldsForGroupDto struct {
	CredentialID        uuid.UUID                `json:"credentialId"`
	GroupID             uuid.UUID                `json:"groupId"`
	UserEncryptedFields []UserEncryptedFieldsDto `json:"userEncryptedFields"`
	AccessType          string                   `json:"accessType"`
}

///////////////////////////////////////////////////////////////////////////////////

type EncryptedCredentialPayload struct {
	CredentialID    uuid.UUID `json:"credentialId" binding:"required"`
	EncryptedFields []Field   `json:"encryptedFields" binding:"required"`
}

type CredentialsForUsersPayload struct {
	UserID         uuid.UUID                    `json:"userId" binding:"required"`
	CredentialData []EncryptedCredentialPayload `json:"credentials" binding:"required"`
	AccessType     string                       `json:"accessType" binding:"required"`
}

type ShareCredentialsWithUsersRequest struct {
	UserData []CredentialsForUsersPayload `json:"userData" binding:"required"`
}

type GroupCredentialPayload struct {
	UserID      uuid.UUID                    `json:"userId" binding:"required"`
	Credentials []EncryptedCredentialPayload `json:"credentials" binding:"required"`
}

type CredentialsForGroupsPayload struct {
	GroupID           uuid.UUID                `json:"groupId" binding:"required"`
	EncryptedUserData []GroupCredentialPayload `json:"encryptedUserData" binding:"required"`
	AccessType        string                   `json:"accessType" binding:"required"`
}

type ShareCredentialsWithGroupsRequest struct {
	GroupData []CredentialsForGroupsPayload `json:"groupData" binding:"required"`
}

type ShareFolderWithUsersRequest struct {
	FolderID uuid.UUID                    `json:"folderId" binding:"required"`
	UserData []CredentialsForUsersPayload `json:"userData" binding:"required"`
}

type ShareFolderWithGroupsRequest struct {
	FolderID  uuid.UUID                     `json:"folderId" binding:"required"`
	GroupData []CredentialsForGroupsPayload `json:"groupData" binding:"required"`
}
