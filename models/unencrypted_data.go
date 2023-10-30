package models

import "github.com/google/uuid"

type UnencryptedData struct {
	BaseModel
	FieldName    string `gorm:"size:255;column:field_name"`
	CredentialID uuid.UUID
	Credential   Credential
	FieldValue   string `gorm:"size:255;column:field_value"`
}

func (u *UnencryptedData) TableName() string {
	return "unencrypted_data"
}
