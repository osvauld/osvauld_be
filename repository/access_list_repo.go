package repository

import (
	"osvauld/infra/database"
	"osvauld/models"
)

func AddAccessList(list *models.AccessList) error {
	db := database.DB
	return db.Create(&list).Error
}
