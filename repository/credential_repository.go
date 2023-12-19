package repository

import (
	"encoding/json"
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/database"
	"osvauld/infra/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetCredentialsByFolder(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) ([]db.FetchCredentialsByUserAndFolderRow, error) {
	arg := db.FetchCredentialsByUserAndFolderParams{
		UserID:   userID,
		FolderID: folderID,
	}
	data, err := database.Store.FetchCredentialsByUserAndFolder(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	return data, nil
}

func ShareCredential(ctx *gin.Context, id uuid.UUID, user dto.User) {
	fieldNames := make([]string, len(user.Fields))
	fieldValues := make([]string, len(user.Fields))
	for idx := range user.Fields {
		fieldNames[idx] = user.Fields[idx].FieldName
		fieldValues[idx] = user.Fields[idx].FieldValue
	}
	params := map[string]interface{}{
		"userId":       user.UserID,
		"credentialId": id,
		"fieldNames":   fieldNames,
		"fieldValues":  fieldValues,
		"accessType":   user.AccessType,
	}
	paramsJson, err := json.Marshal(params)
	if err != nil {
		logger.Errorf(err.Error())
		return
	}
	err = database.Store.ShareSecret(ctx, paramsJson)
	if err != nil {
		logger.Errorf(err.Error())
		return
	}
}

func FetchCredentialByID(ctx *gin.Context, credentialID uuid.UUID) (db.GetCredentialDetailsRow, error) {
	credential, err := database.Store.GetCredentialDetails(ctx, credentialID)
	if err != nil {
		logger.Errorf(err.Error())
		return credential, err
	}
	return credential, nil

}

func FetchEncryptedData(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) ([]db.GetUserEncryptedDataRow, error) {
	arg := db.GetUserEncryptedDataParams{
		CredentialID: credentialID,
		UserID:       userID,
	}
	encryptedData, err := database.Store.GetUserEncryptedData(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	return encryptedData, err
}

func FetchUnEncryptedData(ctx *gin.Context, credentialID uuid.UUID) ([]db.GetCredentialUnencryptedDataRow, error) {

	encryptedData, err := database.Store.GetCredentialUnencryptedData(ctx, credentialID)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	return encryptedData, err
}

func GetEncryptedCredentails(ctx *gin.Context, folderId uuid.UUID, userID uuid.UUID) ([]db.GetEncryptedCredentialsByFolderRow, error) {
	arg := db.GetEncryptedCredentialsByFolderParams{
		FolderID: folderId,
		UserID:   userID,
	}
	encryptedData, err := database.Store.GetEncryptedCredentialsByFolder(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	return encryptedData, err
}

func GetEncryptedCredentailsByIds(ctx *gin.Context, credentialIds []uuid.UUID, userID uuid.UUID) ([]db.GetEncryptedDataByCredentialIdsRow, error) {
	arg := db.GetEncryptedDataByCredentialIdsParams{
		Column1: credentialIds,
		UserID:  userID,
	}
	encryptedData, err := database.Store.GetEncryptedDataByCredentialIds(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	return encryptedData, err
}
