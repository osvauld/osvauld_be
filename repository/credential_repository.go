package repository

import (
	"osvauld/infra/database"
	"osvauld/models"
)

func SaveCredential(credential *models.Credential) error {
	db := database.DB
	return db.Create(&credential).Error
}

func SaveEncryptedData(data *models.EncryptedData) error {
	db := database.DB
	return db.Create(&data).Error
}

func SaveUnencryptedData(data *models.UnencryptedData) error {
	db := database.DB
	return db.Create(&data).Error
}
