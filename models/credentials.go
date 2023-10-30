package models

import "github.com/google/uuid"

type Credential struct {
	BaseModel
	Name        string `gorm:"size:255;column:name"`
	Description string `gorm:"size:255;column:description"`
	FolderID    uuid.UUID
	Folder      Folder
	CreatedBy   uuid.UUID
	Creator     User `gorm:"foreignKey:CreatedBy"`
}

func (u *Credential) TableName() string {
	return "credentials"
}
