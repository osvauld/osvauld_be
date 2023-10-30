package models

import "github.com/google/uuid"

type AccessList struct {
	BaseModel
	CredentialID uuid.UUID
	Credential   Credential
	UserID       uuid.UUID
	User         User
	AccessType   string `gorm:"size:255;column:access_type"`
}

func (u *AccessList) TableName() string {
	return "access_list"
}
