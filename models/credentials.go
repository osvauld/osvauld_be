package models

import "github.com/google/uuid"

type Credential struct {
	BaseModel
	Name             string `gorm:"size:255;column:name"`
	Description      string `gorm:"size:255;column:description"`
	FolderID         uuid.UUID
	Folder           Folder
	CreatedBy        uuid.UUID
	Creator          User              `gorm:"foreignKey:CreatedBy"`
	EncryptedDatas   []EncryptedData   `gorm:"foreignKey:CredentialID"`
	UnencryptedDatas []UnencryptedData `gorm:"foreignKey:CredentialID"`
}

func (u *Credential) TableName() string {
	return "credentials"
}
