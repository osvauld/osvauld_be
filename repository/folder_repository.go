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
	id, err := database.Q.CreateFolder(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return uuid.Nil, err
	}
	return id, nil
}

func GetAccessibleFolders(ctx *gin.Context, userID uuid.UUID) ([]db.FetchAccessibleAndCreatedFoldersByUserRow, error) {
	folders, err := database.Q.FetchAccessibleAndCreatedFoldersByUser(ctx, uuid.NullUUID{UUID: userID, Valid: true})
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	return folders, nil
}

func GetUsersByFolder(ctx *gin.Context, folderID uuid.UUID) ([]db.GetUsersByFolderRow, error) {
	users, err := database.Q.GetUsersByFolder(ctx, uuid.NullUUID{UUID: folderID, Valid: true})
	if err != nil {
		logger.Errorf(err.Error())
		return users, err
	}
	return users, nil
}
