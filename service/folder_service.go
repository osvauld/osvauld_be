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
	"unauthorized": 0,
	"reader":       1,
	"manager":      99,
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

func GetSharedUsersForFolder(ctx *gin.Context, folderID uuid.UUID, caller uuid.UUID) ([]db.GetSharedUsersForFolderRow, error) {

	if err := VerifyFolderReadAccessForUser(ctx, folderID, caller); err != nil {
		return nil, err
	}

	users, err := repository.GetSharedUsersForFolder(ctx, folderID)
	return users, err
}

func GetSharedGroupsForFolder(ctx *gin.Context, folderID uuid.UUID, caller uuid.UUID) ([]db.GetSharedGroupsForFolderRow, error) {

	if err := VerifyFolderReadAccessForUser(ctx, folderID, caller); err != nil {
		return nil, err
	}

	groups, err := repository.GetSharedGroupsForFolder(ctx, folderID)
	return groups, err
}

func GetFolderAccessForUser(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) (string, error) {
	accessValues, err := repository.GetFolderAccessForUser(ctx, db.GetFolderAccessForUserParams{
		FolderID: folderID,
		UserID:   userID,
	})
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

func GetGroupsWithoutAccess(ctx *gin.Context, folderID uuid.UUID, caller uuid.UUID) ([]db.GetGroupsWithoutAccessRow, error) {

	if err := VerifyFolderReadAccessForUser(ctx, folderID, caller); err != nil {
		return nil, err
	}

	groups, err := repository.GetGroupsWithoutAccess(ctx, folderID)
	return groups, err
}
