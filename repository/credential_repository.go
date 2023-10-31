package repository

import (
	"osvauld/infra/database"
	"osvauld/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

func GetFolderIDsByCredentialIDs(credentialIDs []uuid.UUID) ([]uuid.UUID, error) {
	var folderIDs []uuid.UUID
	db := database.DB
	err := db.Table("credentials").
		Where("id IN (?)", credentialIDs).
		Pluck("DISTINCT(folder_id)", &folderIDs).
		Error
	return folderIDs, err
}

func FetchSecretByID(credentialID string) (CustomOutput, error) {
	db := database.DB
	var credential models.Credential

	err := db.
		Select("ID, Name, Description"). // Explicitly selecting fields at the root
		Preload("EncryptedDatas", func(db *gorm.DB) *gorm.DB {
			return db.Omit("Credential", "Folder")
		}).
		Preload("UnencryptedDatas", func(db *gorm.DB) *gorm.DB {
			return db.Omit("Credential", "Folder")
		}).
		Where("id = ?", credentialID).First(&credential).Error
	returnValue := ConvertToCustomOutput(credential)
	return returnValue, err
}

type CustomOutput struct {
	Name              string
	Description       string
	ID                uuid.UUID
	EncryptedFields   map[string]string
	UnencryptedFields map[string]string
}

func ConvertToCustomOutput(cred models.Credential) CustomOutput {
	output := CustomOutput{
		Name:        cred.Name,
		Description: cred.Description,
		ID:          cred.ID,
	}

	// Assign encrypted fields
	output.EncryptedFields = make(map[string]string)
	for _, data := range cred.EncryptedDatas {
		output.EncryptedFields[data.FieldName] = data.FieldValue
	}

	// Assign unencrypted fields
	output.UnencryptedFields = make(map[string]string)
	for _, data := range cred.UnencryptedDatas {
		output.UnencryptedFields[data.FieldName] = data.FieldValue
	}

	return output
}
