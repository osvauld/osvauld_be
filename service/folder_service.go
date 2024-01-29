package service

import (
	"database/sql"
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var folderAccessLevels = map[string]int{
	"member": 0,
	"owner":  99,
}

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

func GetUsersByFolder(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) ([]db.GetUsersByFolderRow, error) {
	users, err := repository.GetUsersByFolder(ctx, folderID)
	return users, err
}

func GetSharedUsersForFolder(ctx *gin.Context, folderID uuid.UUID) ([]db.GetSharedUsersForFolderRow, error) {
	users, err := repository.GetSharedUsersForFolder(ctx, folderID)
	return users, err
}

func GetSharedGroupsForFolder(ctx *gin.Context, folderID uuid.UUID) ([]db.GetSharedGroupsForFolderRow, error) {
	groups, err := repository.GetSharedGroupsForFolder(ctx, folderID)
	return groups, err
}

func GetFolderAccessForUser(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) (string, error) {
	accessValues, err := repository.GetFolderAccessForUser(ctx, folderID, userID)
	if err != nil {
		return "", err
	}

	// Find the higest access level for a user
	highestAccess := "unauthorized"
	for _, accessValue := range accessValues {

		if folderAccessLevels[accessValue] > folderAccessLevels[highestAccess] {
			highestAccess = accessValue
		}
	}

	return highestAccess, nil
}

func CheckFolderOwner(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) (bool, error) {
	access, err := GetFolderAccessForUser(ctx, folderID, userID)
	if err != nil {
		return false, err
	}

	if access == "owner" {
		return true, nil
	}

	return false, nil
}

func CheckOwnerOrManagerAccessForFolder(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) (bool, error) {
	return repository.CheckOwnerOrManagerAccessForFolder(ctx, folderID, userID)

}

func GetGroupsWithoutAccess(ctx *gin.Context, folderID uuid.UUID) ([]db.GetGroupsWithoutAccessRow, error) {
	groups, err := repository.GetGroupsWithoutAccess(ctx, folderID)
	return groups, err
}
