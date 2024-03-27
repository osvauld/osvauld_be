package repository

import (
	db "osvauld/db/sqlc"
	"osvauld/infra/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func HasManageAccessForFolder(ctx *gin.Context, args db.HasManageAccessForFolderParams) (bool, error) {
	return database.Store.HasManageAccessForFolder(ctx, args)
}

func HasReadAccessForFolder(ctx *gin.Context, args db.HasReadAccessForFolderParams) (bool, error) {
	return database.Store.HasReadAccessForFolder(ctx, args)
}

func GetFolderAccess(ctx *gin.Context, folderId uuid.UUID) ([]db.GetAccessTypeAndUserByFolderRow, error) {
	return database.Store.GetAccessTypeAndUserByFolder(ctx, folderId)
}

func GetFolderAccessForUser(ctx *gin.Context, args db.GetFolderAccessForUserParams) ([]string, error) {
	return database.Store.GetFolderAccessForUser(ctx, args)
}

func CheckFolderAccessEntryExists(ctx *gin.Context, args db.CheckFolderAccessEntryExistsParams) (bool, error) {
	return database.Store.CheckFolderAccessEntryExists(ctx, args)
}

func GetFolderUsersWithDirectAccess(ctx *gin.Context, folderId uuid.UUID) ([]db.GetFolderUsersWithDirectAccessRow, error) {
	return database.Store.GetFolderUsersWithDirectAccess(ctx, folderId)
}

func GetFolderGroups(ctx *gin.Context, folderId uuid.UUID) ([]db.GetFolderGroupsRow, error) {
	return database.Store.GetFolderGroups(ctx, folderId)
}

func GetFolderUsersForDataSync(ctx *gin.Context, folderId uuid.UUID) ([]db.GetFolderUsersForDataSyncRow, error) {
	return database.Store.GetFolderUsersForDataSync(ctx, folderId)
}
