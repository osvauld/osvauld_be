package models

import (
	"github.com/google/uuid"
)

type Folder struct {
	BaseModel
	Name        string `gorm:"size:255;column:name"`
	Description string `gorm:"size:255;column:description"`
	CreatedBy   uuid.UUID
	Creator     User `gorm:"foreignKey:CreatedBy"`
}

func (u *Folder) TableName() string {
	return "folders"
}
