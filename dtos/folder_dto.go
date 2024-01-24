package dto

import "github.com/google/uuid"

type CreateFolder struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ShareFolder struct {
	FolderID uuid.UUID    `json:"folderId"`
	Users    []AccessType `json:"users"`
}

type AccessType struct {
	UserID     uuid.UUID `json:"userId"`
	AccessType string    `json:"accessType"`
}

type UserFolderAccessDto struct {
	FolderID   uuid.UUID     `json:"folderId"`
	AccessType string        `json:"accessType"`
	UserID     uuid.UUID     `json:"userId"`
	GroupID    uuid.NullUUID `json:"groupId"`
}
