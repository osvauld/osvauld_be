package service

import (
	"database/sql"
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var accessLevels = map[string]int{
	"manager": 2,
	"reader":  1,
	"none":    0,
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

func FetchAccessibleFoldersForUser(ctx *gin.Context, userID uuid.UUID) ([]dto.FolderDetails, error) {
	folders, err := repository.FetchAccessibleFoldersForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	uniqueFolders, err := GetFolderWithHighestAccess(folders)
	if err != nil {
		return nil, err
	}

	uniqueFoldersDtos := []dto.FolderDetails{}
	for _, folder := range uniqueFolders {
		uniqueFoldersDtos = append(uniqueFoldersDtos, dto.FolderDetails{
			FolderID:    folder.ID,
			Name:        folder.Name,
			Description: folder.Description.String,
			CreatedAt:   folder.CreatedAt,
			CreatedBy:   folder.CreatedBy.UUID,
			AccessType:  folder.AccessType,
		})
	}

	return uniqueFoldersDtos, nil
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

func EditFolder(ctx *gin.Context, folderID uuid.UUID, payload dto.EditFolder, caller uuid.UUID) error {

	if err := VerifyFolderManageAccessForUser(ctx, folderID, caller); err != nil {
		return err
	}

	args := db.EditFolderParams{
		ID:   folderID,
		Name: payload.Name,
		Description: sql.NullString{
			String: payload.Description,
			Valid:  true,
		},
	}

	err := repository.EditFolder(ctx, args)
	if err != nil {
		return err
	}
	return nil
}

func GetFolderWithHighestAccess(folders []db.FetchAccessibleFoldersForUserRow) ([]db.FetchAccessibleFoldersForUserRow, error) {
	uniqueFoldersMap := make(map[uuid.UUID]db.FetchAccessibleFoldersForUserRow)

	for _, folder := range folders {
		if existingFolder, ok := uniqueFoldersMap[folder.ID]; ok {
			if accessLevels[folder.AccessType] > accessLevels[existingFolder.AccessType] {
				uniqueFoldersMap[folder.ID] = folder
			}
		} else {
			uniqueFoldersMap[folder.ID] = folder
		}
	}
	uniqueFolders := make([]db.FetchAccessibleFoldersForUserRow, 0, len(uniqueFoldersMap))
	for _, folder := range uniqueFoldersMap {
		uniqueFolders = append(uniqueFolders, folder)
	}

	return uniqueFolders, nil

}
