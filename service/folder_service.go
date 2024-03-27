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

func RemoveFolder(ctx *gin.Context, folderID uuid.UUID, caller uuid.UUID) error {
	// TODO: check folder ownership before removal
	if err := VerifyFolderManageAccessForUser(ctx, folderID, caller); err != nil {
		return err
	}
	err := repository.RemoveFolder(ctx, folderID)
	if err != nil {
		return err
	}
	return nil
}
