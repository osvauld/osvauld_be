package dto

import "github.com/google/uuid"

type RemoveCredentialAccessForUsers struct {
	UserIDs []uuid.UUID `json:"userIds"`
}

type RemoveFolderAccessForUsers struct {
	UserIDs []uuid.UUID `json:"userIds"`
}

type RemoveCredentialAccessForGroups struct {
	GroupIDs []uuid.UUID `json:"groupIds"`
}

type RemoveFolderAccessForGroups struct {
	GroupIDs []uuid.UUID `json:"groupIds"`
}

type EditCredentialAccessForUser struct {
	UserID     uuid.UUID `json:"userId"`
	AccessType string    `json:"accessType"`
}

type EditFolderAccessForUser struct {
	UserID     uuid.UUID `json:"userId"`
	AccessType string    `json:"accessType"`
}

type EditCredentialAccessForGroup struct {
	GroupID    uuid.UUID `json:"groupId"`
	AccessType string    `json:"accessType"`
}

type EditFolderAccessForGroup struct {
	GroupID    uuid.UUID `json:"groupId"`
	AccessType string    `json:"accessType"`
}

type CredentialUserWithAccess struct {
	UserID       uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	AccessType   string    `json:"accessType"`
	AccessSource string    `json:"accessSource"`
}

type FolderUserWithAccess struct {
	UserID       uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	AccessType   string    `json:"accessType"`
	AccessSource string    `json:"accessSource"`
}
