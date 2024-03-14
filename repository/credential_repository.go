package repository

import (
	db "osvauld/db/sqlc"
	"osvauld/infra/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddCredential(ctx *gin.Context, args db.AddCredentialTransactionParams) (uuid.UUID, error) {

	return database.Store.AddCredentialTransaction(ctx, args)

}

func GetCredentialDataByID(ctx *gin.Context, credentialID uuid.UUID) (db.GetCredentialDataByIDRow, error) {

	return database.Store.GetCredentialDataByID(ctx, credentialID)
}

func FetchCredentialDetailsForUserByFolderId(ctx *gin.Context, args db.FetchCredentialDetailsForUserByFolderIdParams) ([]db.FetchCredentialDetailsForUserByFolderIdRow, error) {

	return database.Store.FetchCredentialDetailsForUserByFolderId(ctx, args)

}

func EditCredential(ctx *gin.Context, args db.EditCredentialTransactionParams) error {

	return database.Store.EditCredentialTransaction(ctx, args)

}

func GetCredentialDetailsByIDs(ctx *gin.Context, credentialIDs []uuid.UUID) ([]db.GetCredentialDetailsByIDsRow, error) {

	return database.Store.GetCredentialDetailsByIDs(ctx, credentialIDs)
}

func GetAllUrlsForUser(ctx *gin.Context, userID uuid.UUID) ([]db.GetAllUrlsForUserRow, error) {

	return database.Store.GetAllUrlsForUser(ctx, userID)

}

func GetCredentialIdsByFolderAndUserId(ctx *gin.Context, args db.GetCredentialIdsByFolderParams) ([]uuid.UUID, error) {

	return database.Store.GetCredentialIdsByFolder(ctx, args)

}

func GetCredentialUsers(ctx *gin.Context, credentialID uuid.UUID) ([]db.GetAccessTypeAndUsersByCredentialIdRow, error) {
	return database.Store.GetAccessTypeAndUsersByCredentialId(ctx, credentialID)

}

func GetCredentialGroups(ctx *gin.Context, credentialID uuid.UUID) ([]db.GetAccessTypeAndGroupsByCredentialIdRow, error) {
	return database.Store.GetAccessTypeAndGroupsByCredentialId(ctx, credentialID)

}

func GetSearchData(ctx *gin.Context, userID uuid.UUID) ([]db.GetCredentialsForSearchByUserIDRow, error) {
	return database.Store.GetCredentialsForSearchByUserID(ctx, userID)
}
