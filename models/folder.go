package models

import (
	"github.com/google/uuid"
)

type Folder struct {
	BaseModel
	Name      string   `gorm:"size:255;column:name"`
	Tags      []string `gorm:"type:varchar(255)[];column:tags"`
	CreatedBy uuid.UUID
	Creator   User `gorm:"foreignKey:CreatedBy"`
}

func (u *Folder) TableName() string {
	return "folders"
}
