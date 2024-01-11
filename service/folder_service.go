package service

import (
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

func GetSharedUsers(ctx *gin.Context, folderID uuid.UUID) ([]db.GetSharedUsersRow, error) {
	users, err := repository.GetSharedUsers(ctx, folderID)
	return users, err
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
