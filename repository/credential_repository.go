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

func AddCredential(ctx *gin.Context, payload dto.SQLCPayload) (uuid.UUID, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		logger.Errorf(err.Error())
		return uuid.Nil, err
	}
	credentialId, err := database.Q.AddCredential(ctx, payloadJSON)
	if err != nil {
		logger.Errorf(err.Error())
		return uuid.Nil, err
	}

	return credentialId, nil
}

func GetCredentialsByFolder(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) ([]db.FetchCredentialsByUserAndFolderRow, error) {
	arg := db.FetchCredentialsByUserAndFolderParams{
		UserID:   uuid.NullUUID{UUID: userID, Valid: true},
		FolderID: uuid.NullUUID{UUID: folderID, Valid: true},
	}
	data, err := database.Q.FetchCredentialsByUserAndFolder(ctx, arg)
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

	arg := db.ShareSecretParams{
		PUserID:       user.UserID,
		PCredentialID: id,
		PFieldNames:   GoSliceToPostgresArray(fieldNames),
		PFieldValues:  GoSliceToPostgresArray(fieldValues),
		PAccessType:   user.AccessType,
	}
	err := database.Q.ShareSecret(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return
	}
}

func FetchCredentialByID(ctx *gin.Context, credentialID uuid.UUID) (db.GetCredentialDetailsRow, error) {
	credential, err := database.Q.GetCredentialDetails(ctx, credentialID)
	if err != nil {
		logger.Errorf(err.Error())
		return credential, err
	}
	return credential, nil

}

func FetchEncryptedData(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) ([]db.GetUserEncryptedDataRow, error) {
	arg := db.GetUserEncryptedDataParams{
		CredentialID: uuid.NullUUID{UUID: credentialID, Valid: true},
		UserID:       uuid.NullUUID{UUID: userID, Valid: true},
	}
	encryptedData, err := database.Q.GetUserEncryptedData(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	return encryptedData, err
}

func FetchUnEncryptedData(ctx *gin.Context, credentialID uuid.UUID) ([]db.GetCredentialUnencryptedDataRow, error) {

	encryptedData, err := database.Q.GetCredentialUnencryptedData(ctx, uuid.NullUUID{UUID: credentialID, Valid: true})
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	return encryptedData, err
}
