package repository

import (
	"fmt"
	"osvauld/infra/database"
	"osvauld/models"

	"github.com/google/uuid"
)

func SaveFolder(folder *models.Folder) error {
	db := database.DB
	fmt.Printf("Folder data before saving %+v\n", folder)
	result := db.Create(&folder)
	fmt.Println(result.Error)
	return result.Error
}

func GetFoldersByIds(folderIds []uuid.UUID) ([]models.Folder, error) {
	db := database.DB
	var folderList []models.Folder
	err := db.Where("id IN (?)", folderIds).Find(&folderList).Error
	return folderList, err
}
