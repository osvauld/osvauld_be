package service

import (
	"database/sql"
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


func CreateFolder(ctx *gin.Context, folder dto.CreateFolderRequest, caller uuid.UUID) (dto.FolderDetails, error) {

	createFolderParams := db.CreateFolderTransactionParams{
		Name:        folder.Name,
		Description: sql.NullString{String: folder.Description, Valid: true},
		CreatedBy:   caller,
	}

	folderDetails, err := repository.CreateFolder(ctx, createFolderParams)
	if err != nil {
		return dto.FolderDetails{}, err
	}
	return folderDetails, nil
}

func FetchAccessibleFoldersForUser(ctx *gin.Context, userID uuid.UUID) ([]db.FetchAccessibleFoldersForUserRow, error) {
	folders, err := repository.FetchAccessibleFoldersForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return folders, nil
}
