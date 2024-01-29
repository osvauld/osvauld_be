package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateFolderRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type FolderDetails struct {
	FolderID    uuid.UUID `json:"folderId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   uuid.UUID `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
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
