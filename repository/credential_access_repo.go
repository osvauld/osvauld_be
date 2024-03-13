package repository

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/database"
	"osvauld/infra/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddCredentialAccess(ctx *gin.Context, credentialID uuid.UUID, accessType string, userID uuid.UUID) error {
	arg := db.AddCredentialAccessParams{
		CredentialID: credentialID,
		UserID:       userID,
	}
	_, err := database.Store.AddCredentialAccess(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	return nil
}

func GetUsersByCredential(ctx *gin.Context, credentailID uuid.UUID) ([]db.GetUsersByCredentialRow, error) {
	users, err := database.Store.GetUsersByCredential(ctx, credentailID)
	if err != nil {
		logger.Errorf(err.Error())
		return users, err
	}
	return users, nil
}

func GetCredentialAccessForUser(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) ([]dto.AccessListResult, error) {
	arg := db.GetCredentialAccessForUserParams{
		CredentialID: credentialID,
		UserID:       userID,
	}

	access, err := database.Store.GetCredentialAccessForUser(ctx, arg)
	if err != nil {
		return []dto.AccessListResult{}, err
	}
	accessListResults := []dto.AccessListResult{}
	for _, access := range access {
		accessListResults = append(accessListResults, dto.AccessListResult{
			ID:         access.ID,
			UserID:     access.UserID,
			AccessType: access.AccessType,
			GroupID:    access.GroupID,
		})
	}
	return accessListResults, nil
}

func CheckCredentialAccessEntryExists(ctx *gin.Context, args db.CheckCredentialAccessEntryExistsParams) (bool, error) {

	return database.Store.CheckCredentialAccessEntryExists(ctx, args)

}

func RemoveCredentialAccessForUsers(ctx *gin.Context, args db.RemoveCredentialAccessForUsersParams) error {

	return database.Store.RemoveCredentialAccessForUsers(ctx, args)
}

func RemoveFolderAccessForUser(ctx *gin.Context, args db.RemoveFolderAccessForUsersParams) error {

	return database.Store.RemoveFolderAccessForUsersTransactions(ctx, args)
}

func RemoveCredentialAccessForGroups(ctx *gin.Context, args db.RemoveCredentialAccessForGroupsParams) error {

	return database.Store.RemoveCredentialAccessForGroups(ctx, args)
}

func RemoveFolderAccessForGroups(ctx *gin.Context, args db.RemoveFolderAccessForGroupsParams) error {

	return database.Store.RemoveFolderAccessForGroupsTransactions(ctx, args)
}

func EditCredentialAccessForUsers(ctx *gin.Context, args db.EditCredentialAccessForUserParams) error {

	return database.Store.EditCredentialAccessForUser(ctx, args)
}

func EditFolderAccessForUser(ctx *gin.Context, args db.EditFolderAccessForUserParams) error {

	return database.Store.EditFolderAccessForUserTransaction(ctx, args)
}

func EditCredentialAccessForGroup(ctx *gin.Context, args db.EditCredentialAccessForGroupParams) error {

	return database.Store.EditCredentialAccessForGroup(ctx, args)
}

func EditFolderAccessForGroup(ctx *gin.Context, args db.EditFolderAccessForGroupParams) error {

	return database.Store.EditFolderAccessForGroupTransaction(ctx, args)
}
