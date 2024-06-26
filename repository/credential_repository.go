package repository

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
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

func EditCredential(ctx *gin.Context, args dto.EditCredentialRequest, caller uuid.UUID) error {
	return database.Store.EditCredentialTransaction(ctx, args, caller)
}

func EditCredentialDetails(ctx *gin.Context, args db.EditCredentialDetailsParams) error {
	return database.Store.EditCredentialDetails(ctx, args)
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

func GetSearchData(ctx *gin.Context, userID uuid.UUID) ([]db.GetCredentialsForSearchByUserIDRow, error) {
	return database.Store.GetCredentialsForSearchByUserID(ctx, userID)
}

func RemoveCredential(ctx *gin.Context, credentialID uuid.UUID) error {
	return database.Store.RemoveCredential(ctx, credentialID)
}
