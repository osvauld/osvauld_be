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
	id, err := database.Store.CreateFolder(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return uuid.Nil, err
	}
	return id, nil
}

func GetAccessibleFolders(ctx *gin.Context, userID uuid.UUID) ([]db.FetchAccessibleAndCreatedFoldersByUserRow, error) {
	folders, err := database.Store.FetchAccessibleAndCreatedFoldersByUser(ctx, uuid.NullUUID{UUID: userID, Valid: true})
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

func ShareFolder(ctx *gin.Context, folder dto.ShareFolder) error {

	var users []uuid.UUID
	var accessTypes []string

	for _, userAccess := range folder.Users {
		users = append(users, userAccess.UserID)
		accessTypes = append(accessTypes, userAccess.AccessType)
	}
	arg := db.AddFolderAccessParams{
		FolderID: folder.FolderID,
		Column2:  users,
		Column3:  accessTypes,
	}
	err := database.Store.AddFolderAccess(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	return nil
}

func GetSharedUsers(ctx *gin.Context, folderID uuid.UUID) ([]db.GetSharedUsersRow, error) {
	users, err := database.Store.GetSharedUsers(ctx, folderID)
	if err != nil {
		logger.Errorf(err.Error())
		return users, err
	}
	return users, nil
}
