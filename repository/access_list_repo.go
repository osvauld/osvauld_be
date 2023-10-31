package repository

import (
	"osvauld/infra/database"
	"osvauld/models"

	"github.com/google/uuid"
)

func AddAccessList(list *models.AccessList) error {
	db := database.DB
	return db.Create(&list).Error
}

func GetCredentialIDsByUserID(user_id uuid.UUID) ([]uuid.UUID, error) {
	db := database.DB
	var credentialIDs []uuid.UUID
	err := db.Table("access_list").Where("user_id = ?", user_id).Pluck("credential_id", &credentialIDs).Error
	return credentialIDs, err
}
