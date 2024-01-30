package dto

import "github.com/google/uuid"

type UserFields struct {
	UserID uuid.UUID `json:"userId"`
	Fields []Field   `json:"fields"`
}

type UserFieldsWithAccessType struct {
	UserID     uuid.UUID `json:"userId"`
	Fields     []Field   `json:"fields"`
	AccessType string    `json:"accessType"`
}

type Field struct {
	ID         uuid.UUID `json:"id"`
	FieldName  string    `json:"fieldName"`
	FieldValue string    `json:"fieldValue"`
	FieldType  string    `json:"fieldType"`
}

type ShareField struct {
	ID         uuid.UUID `json:"fieldId"`
	FieldValue string    `json:"fieldValue"`
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

type ShareCredentialsWithUsersRequest struct {
	UserData []ShareCredentialsForUserPayload `json:"userData" binding:"required"`
}

type CredentialsForGroupUsersPayload struct {
	UserID         uuid.UUID                `json:"userId" binding:"required"`
	CredentialData []ShareCredentialPayload `json:"credentials" binding:"required"`
}

type CredentialsForGroupsPayload struct {
	GroupID    uuid.UUID                        `json:"groupId" binding:"required"`
	AccessType string                           `json:"accessType" binding:"required"`
	UserData   []ShareCredentialsForUserPayload `json:"userData" binding:"required"`
}

type ShareCredentialsWithGroupsRequest struct {
	GroupData []CredentialsForGroupsPayload `json:"groupData" binding:"required"`
}

type GroupCredentialPayload struct {
	UserID      uuid.UUID                    `json:"userId" binding:"required"`
	Credentials []EncryptedCredentialPayload `json:"credentials" binding:"required"`
}

type CredentialFieldsForUserDto struct {
	CredentialID uuid.UUID     `json:"credentialId"`
	UserID       uuid.UUID     `json:"userId"`
	Fields       []Field       `json:"Fields"`
	AccessType   string        `json:"accessType"`
	GroupID      uuid.NullUUID `json:"groupId"`
	FolderID     uuid.NullUUID `json:"folderId"`
}

type CredentialFieldsForGroupDto struct {
	CredentialID        uuid.UUID    `json:"credentialId"`
	GroupID             uuid.UUID    `json:"groupId"`
	UserEncryptedFields []UserFields `json:"userFields"`
	AccessType          string       `json:"accessType"`
}

type ShareFolderForUserPayload struct {
	UserID         uuid.UUID                `json:"userId" binding:"required"`
	CredentialData []ShareCredentialPayload `json:"credentials" binding:"required"`
}

type ShareFolderWithUsersRequest struct {
	FolderID uuid.UUID                        `json:"folderId" binding:"required"`
	UserData []ShareCredentialsForUserPayload `json:"userData" binding:"required"`
}

type ShareFolderWithGroupsRequest struct {
	FolderID  uuid.UUID                     `json:"folderId" binding:"required"`
	GroupData []CredentialsForGroupsPayload `json:"groupData" binding:"required"`
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
