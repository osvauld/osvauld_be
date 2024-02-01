package repository

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/database"
	"osvauld/infra/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddCredential(ctx *gin.Context, args dto.AddCredentialDto, caller uuid.UUID) (uuid.UUID, error) {

	credentialID, err := database.Store.AddCredentialTransaction(ctx, args, caller)
	if err != nil {
		return uuid.UUID{}, err
	}
	return credentialID, nil
}

func FetchCredentialByID(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) (dto.CredentialDetails, error) {

	credentialDetails := dto.CredentialDetails{}
	credentialDetails.UserID = userID

	credential, err := database.Store.FetchCredentialDataByID(ctx, credentialID)
	if err != nil {
		return dto.CredentialDetails{}, err
	}

	credentialDetails.CredentialID = credential.ID
	credentialDetails.Name = credential.Name
	credentialDetails.FolderID = credential.FolderID
	credentialDetails.CreatedAt = credential.CreatedAt
	credentialDetails.UpdatedAt = credential.UpdatedAt
	credentialDetails.CreatedBy = credential.CreatedBy
	if credential.Description.Valid {
		credentialDetails.Description = credential.Description.String
	}

	return credentialDetails, nil
}

func FetchCredentialDetailsForUserByFolderId(ctx *gin.Context, args db.FetchCredentialDetailsForUserByFolderIdParams) ([]db.FetchCredentialDetailsForUserByFolderIdRow, error) {

	return database.Store.FetchCredentialDetailsForUserByFolderId(ctx, args)

}

func FetchUnEncryptedData(ctx *gin.Context, credentialID uuid.UUID) ([]db.GetCredentialUnencryptedDataRow, error) {

	encryptedData, err := database.Store.GetCredentialUnencryptedData(ctx, credentialID)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	return encryptedData, err
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

func GetSensitiveFieldsById(ctx *gin.Context, credentialID uuid.UUID, caller uuid.UUID) ([]db.GetSensitiveFieldsRow, error) {
	// Check if caller has access
	sensitiveFields, err := database.Store.GetSensitiveFields(ctx, db.GetSensitiveFieldsParams{
		CredentialID: credentialID,
		UserID:       caller,
	})

	return sensitiveFields, err
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
