package models

type User struct {
	BaseModel
	Username string `gorm:"size:255;unique;column:username"`
}

func (u *User) TableName() string {
	return "users"
}
