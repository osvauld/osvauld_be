package migrations

import (
	"osvauld/infra/database"
	"osvauld/infra/logger"
	"osvauld/models"
)

// Migrate Add list of model add for migrations
// TODO later separate migration each models
func Migrate() {
	var migrationModels = []interface{}{&models.Example{}, &models.UnencryptedData{}, &models.AccessList{}, &models.Credential{}, &models.EncryptedData{}, &models.Folder{}, &models.User{}}
	database.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	err := database.DB.AutoMigrate(migrationModels...)
	if err != nil {
		logger.Errorf("Migration Error %v", err)
		return
	}
}
