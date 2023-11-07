package service

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateFolder(ctx *gin.Context, folder dto.CreateFolder, userID uuid.UUID) error {
	_, err := repository.CreateFolder(ctx, folder, userID)
	if err != nil {
		return err
	}
	return nil
}

func GetAccessibleFolders(ctx *gin.Context, userID uuid.UUID) ([]db.FetchAccessibleAndCreatedFoldersByUserRow, error) {
	folders, err := repository.GetAccessibleFolders(ctx, userID)
	if err != nil {
		return nil, err
	}
	return folders, nil
}

func GetUsersByFolder(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) ([]db.GetUsersByFolderRow, error) {
	users, err := repository.GetUsersByFolder(ctx, folderID)
	return users, err
}
