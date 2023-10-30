package models

import "github.com/google/uuid"

type EncryptedData struct {
	BaseModel
	FieldName    string `gorm:"size:255;column:field_name"`
	CredentialID uuid.UUID
	Credential   Credential
	FieldValue   string `gorm:"size:255;column:field_value"`
	UserID       uuid.UUID
	User         User
}

func (u *EncryptedData) TableName() string {
	return "encrypted_data"
}
