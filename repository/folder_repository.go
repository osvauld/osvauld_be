package repository

import (
	"fmt"
	"osvauld/infra/database"
	"osvauld/models"
)

// SaveUser saves a user to the database.
func SaveFolder(folder *models.Folder) error {
	db := database.DB
	fmt.Printf("Folder data before saving %+v\n", folder)
	result := db.Create(&folder)
	fmt.Println(result.Error)
	return result.Error
}

// Other CRUD operations (GetUser, UpdateUser, DeleteUser, etc.) can be added here.
