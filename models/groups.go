package models

import "github.com/google/uuid"

type Group struct {
	BaseModel
	Name    string      `gorm:"size:255;column:name"`
	Members []uuid.UUID `gorm:"type:uuid[];column:members"`
}

func (u *Group) TableName() string {
	return "groups"
}
