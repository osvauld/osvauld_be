package repository

import (
	db "osvauld/db/sqlc"
	"osvauld/infra/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUsersByCredential(ctx *gin.Context, credentailID uuid.UUID) ([]db.GetUsersByCredentialRow, error) {

	return database.Store.GetUsersByCredential(ctx, credentailID)
}

func CheckCredentialAccessEntryExists(ctx *gin.Context, args db.CheckCredentialAccessEntryExistsParams) (bool, error) {

	return database.Store.CheckCredentialAccessEntryExists(ctx, args)

}

func GetCredentialAccessTypeForUser(ctx *gin.Context, args db.GetCredentialAccessTypeForUserParams) ([]db.GetCredentialAccessTypeForUserRow, error) {

	return database.Store.GetCredentialAccessTypeForUser(ctx, args)
}

func HasManageAccessForCredential(ctx *gin.Context, args db.HasManageAccessForCredentialParams) (bool, error) {

	return database.Store.HasManageAccessForCredential(ctx, args)
}

func HasReadAccessForCredential(ctx *gin.Context, args db.HasReadAccessForCredentialParams) (bool, error) {

	return database.Store.HasReadAccessForCredential(ctx, args)
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

func GetCredentialUsersWithDirectAccess(ctx *gin.Context, credentialID uuid.UUID) ([]db.GetCredentialUsersWithDirectAccessRow, error) {

	return database.Store.GetCredentialUsersWithDirectAccess(ctx, credentialID)
}

func GetCredentialUsersWithDirectAndGroupAccess(ctx *gin.Context, credentialID uuid.UUID) ([]db.GetCredentialUsersWithDirectAndGroupAccessRow, error) {

	return database.Store.GetCredentialUsersWithDirectAndGroupAccess(ctx, credentialID)
}
