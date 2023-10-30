package repository

import (
	"osvauld/infra/database"
	"osvauld/models"
)

// SaveUser saves a user to the database.
func SaveUser(user *models.User) error {
	db := database.DB
	result := db.Create(&user)
	return result.Error
}

// Other CRUD operations (GetUser, UpdateUser, DeleteUser, etc.) can be added here.
