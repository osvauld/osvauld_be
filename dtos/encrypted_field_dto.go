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

type UserCredentialPayload struct {
	CredentialID  uuid.UUID `json:"credentialId" binding:"required"`
	EncryptedData []Field   `json:"encryptedData" binding:"required"`
}

type MulitpleCredentialsForUserPayload struct {
	UserID         uuid.UUID               `json:"userId" binding:"required"`
	CredentialData []UserCredentialPayload `json:"credentials" binding:"required"`
	AccessType     string                  `json:"accessType" binding:"required"`
}

type ShareMultipleCredentialsWithMultipleUserRequest struct {
	UserData []MulitpleCredentialsForUserPayload `json:"userData" binding:"required"`
}

type GroupCredentialPayload struct {
	CredentialID uuid.UUID                `json:"credentialId" binding:"required"`
	UserData     []UserEncryptedFieldsDto `json:"userData" binding:"required"`
}

type MulitpleCredentialsForGroupPayload struct {
	GroupID        uuid.UUID                `json:"groupId" binding:"required"`
	CredentialData []GroupCredentialPayload `json:"credentials" binding:"required"`
	AccessType     string                   `json:"accessType" binding:"required"`
}

type ShareMultipleCredentialsWithMultipleGroupRequest struct {
	GroupData []MulitpleCredentialsForGroupPayload `json:"groupData" binding:"required"`
}
