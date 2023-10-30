package repository

import (
	"osvauld/infra/database"
	"osvauld/models"

	"github.com/google/uuid"
)

type SecretOutput struct {
	ID     uuid.UUID                `json:"ID"`
	Fields []models.UnencryptedData `json:"Fields"`
}

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

func GetSecretsByFolderAndUser(folderID, userID uuid.UUID) ([]SecretOutput, error) {
	var results []SecretOutput

	db := database.DB

	// Fetch the credentials for the folder
	var credentials []models.Credential
	db.Where("folder_id = ?", folderID).Find(&credentials)

	// For each credential, check access and fetch secrets if permitted
	for _, credential := range credentials {
		var accessList models.AccessList
		if err := db.Where("credential_id = ? AND user_id = ?", credential.ID, userID).First(&accessList).Error; err == nil {
			// If access found, fetch the secrets
			var secretData []models.UnencryptedData
			db.Where("credential_id = ?", credential.ID).Find(&secretData)

			result := SecretOutput{
				ID:     credential.ID,
				Fields: secretData,
			}
			results = append(results, result)
		}
	}

	return results, nil
}
