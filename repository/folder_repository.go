package repository

import (
	"database/sql"
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/database"
	"osvauld/infra/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateFolder(ctx *gin.Context, folder dto.CreateFolder, userID uuid.UUID) (uuid.UUID, error) {
	arg := db.CreateFolderParams{
		Name:        folder.Name,
		Description: sql.NullString{String: folder.Description, Valid: true},
		CreatedBy:   uuid.NullUUID{UUID: userID, Valid: true},
	}
	q := db.New(database.DB)
	id, err := q.CreateFolder(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return uuid.Nil, err
	}
	return id, nil
}

func GetAccessibleFolders(ctx *gin.Context, userID uuid.UUID) ([]db.Folder, error) {
	q := db.New(database.DB)
	folders, err := q.FetchAccessibleAndCreatedFoldersByUser(ctx, uuid.NullUUID{UUID: userID, Valid: true})
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	return folders, nil
}
