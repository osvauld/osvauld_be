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
