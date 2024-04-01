package repository

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateFolder(ctx *gin.Context, args db.CreateFolderTransactionParams) (dto.FolderDetails, error) {
	return database.Store.CreateFolderTransaction(ctx, args)
}

func FetchAccessibleFoldersForUser(ctx *gin.Context, userID uuid.UUID) ([]db.FetchAccessibleFoldersForUserRow, error) {
	return database.Store.FetchAccessibleFoldersForUser(ctx, userID)
}

func GetGroupsWithoutAccess(ctx *gin.Context, folderId uuid.UUID, caller uuid.UUID) ([]db.GetGroupsWithoutAccessRow, error) {
	return database.Store.GetGroupsWithoutAccess(ctx, db.GetGroupsWithoutAccessParams{
		FolderID: folderId,
		ID:       caller,
	})
}

func RemoveFolder(ctx *gin.Context, folderID uuid.UUID) error {
	return database.Store.RemoveFolder(ctx, folderID)
}
