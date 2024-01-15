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
		CreatedBy:   userID,
	}
	id, err := database.Store.CreateFolder(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return uuid.Nil, err
	}
	return id, nil
}

func GetAccessibleFolders(ctx *gin.Context, userID uuid.UUID) ([]db.FetchAccessibleAndCreatedFoldersByUserRow, error) {
	folders, err := database.Store.FetchAccessibleAndCreatedFoldersByUser(ctx, userID)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	return folders, nil
}

func GetUsersByFolder(ctx *gin.Context, folderID uuid.UUID) ([]db.GetUsersByFolderRow, error) {
	users, err := database.Store.GetUsersByFolder(ctx, folderID)
	if err != nil {
		logger.Errorf(err.Error())
		return users, err
	}
	return users, nil
}

func CheckFolderAccess(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) (bool, error) {

	arg := db.IsFolderOwnerParams{
		UserID:   userID,
		FolderID: folderID,
	}
	access, err := database.Store.IsFolderOwner(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return false, err
	}
	return access, nil
}

func GetSharedUsers(ctx *gin.Context, folderID uuid.UUID) ([]db.GetSharedUsersRow, error) {
	users, err := database.Store.GetSharedUsers(ctx, folderID)
	if err != nil {
		logger.Errorf(err.Error())
		return users, err
	}
	return users, nil
}

func CheckOwnerOrManagerAccessForFolder(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) (bool, error) {
	arg := db.IsUserManagerOrOwnerParams{
		UserID:   userID,
		FolderID: folderID,
	}
	access, err := database.Store.IsUserManagerOrOwner(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return false, err
	}
	return access, nil
}

func GetFolderAccessForUser(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) ([]string, error) {

	params := db.GetFolderAccessForUserParams{
		FolderID: folderID,
		UserID:   userID,
	}
	accessRows, err := database.Store.GetFolderAccessForUser(ctx, params)
	if err != nil {
		logger.Errorf(err.Error())
		return []string{}, err
	}
	return accessRows, nil
}

func GetFolderAccess(ctx *gin.Context, folderId uuid.UUID) ([]db.GetAccessTypeAndUserByFolderRow, error) {
	access, err := database.Store.GetAccessTypeAndUserByFolder(ctx, folderId)
	if err != nil {
		logger.Errorf(err.Error())
		return access, err
	}
	return access, nil
}

func GetGroupsWithoutAccess(ctx *gin.Context, folderId uuid.UUID) ([]db.GetGroupsWithoutAccessRow, error) {
	groups, err := database.Store.GetGroupsWithoutAccess(ctx, folderId)
	if err != nil {
		logger.Errorf(err.Error())
		return groups, err
	}
	return groups, nil
}
