package repository

import (
	db "osvauld/db/sqlc"
	"osvauld/infra/database"
	"osvauld/infra/logger"

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

func GetCredentialsFieldsByIds(ctx *gin.Context, credentialIds []uuid.UUID, userID uuid.UUID) ([]db.GetCredentialsFieldsByIdsRow, error) {
	arg := db.GetCredentialsFieldsByIdsParams{
		Column1: credentialIds,
		UserID:  userID,
	}
	encryptedData, err := database.Store.GetCredentialsFieldsByIds(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	return encryptedData, err
}

func GetCredentialsByIDs(ctx *gin.Context, credentialIds []uuid.UUID, userID uuid.UUID) ([]db.GetCredentialDetailsByIdsRow, error) {
	credentials, err := database.Store.GetCredentialDetailsByIds(ctx, db.GetCredentialDetailsByIdsParams{
		UserID:  userID,
		Column1: credentialIds,
	})
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	return credentials, err
}

func GetAllUrlsForUser(ctx *gin.Context, userID uuid.UUID) ([]db.GetAllUrlsForUserRow, error) {
	urls, err := database.Store.GetAllUrlsForUser(ctx, userID)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	return urls, err
}

func GetCredentialIdsByFolderAndUserId(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) ([]uuid.UUID, error) {
	credentialIds, err := database.Store.GetCredentialIdsByFolder(ctx, db.GetCredentialIdsByFolderParams{
		FolderID: folderID,
		UserID:   userID,
	})
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	return credentialIds, err
}

func GetCredentialUsers(ctx *gin.Context, credentialID uuid.UUID) ([]db.GetAccessTypeAndUsersByCredentialIdRow, error) {
	return database.Store.GetAccessTypeAndUsersByCredentialId(ctx, credentialID)

}

func GetCredentialGroups(ctx *gin.Context, credentialID uuid.UUID) ([]db.GetAccessTypeAndGroupsByCredentialIdRow, error) {
	return database.Store.GetAccessTypeAndGroupsByCredentialId(ctx, credentialID)

}
