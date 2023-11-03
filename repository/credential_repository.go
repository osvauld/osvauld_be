package repository

import (
	"encoding/json"
	"fmt"
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/database"
	"osvauld/infra/logger"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddCredential(ctx *gin.Context, payload dto.SQLCPayload) (uuid.UUID, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		logger.Errorf(err.Error())
		return uuid.Nil, err
	}
	credentialIdInterface, err := database.Q.AddCredential(ctx, payloadJSON)
	if err != nil {
		logger.Errorf(err.Error())
		return uuid.Nil, err
	}
	credentialId, ok := credentialIdInterface.(uuid.UUID)
	if !ok {
		logger.Errorf("Type assertion failed: Expected uuid.UUID, got %T", credentialIdInterface)
		return uuid.Nil, fmt.Errorf("incorrect type returned from AddCredential")
	}

	return credentialId, nil
}

func GetCredentialsByFolder(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) ([]db.FetchCredentialsByUserAndFolderRow, error) {
	arg := db.FetchCredentialsByUserAndFolderParams{
		UserID:   uuid.NullUUID{UUID: userID, Valid: true},
		FolderID: uuid.NullUUID{UUID: folderID, Valid: true},
	}
	q := db.New(database.DB)
	data, err := q.FetchCredentialsByUserAndFolder(ctx, arg)
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
		ShareSecret:   user.UserID,
		ShareSecret_2: id,
		ShareSecret_3: GoSliceToPostgresArray(fieldNames),
		ShareSecret_4: GoSliceToPostgresArray(fieldValues),
		ShareSecret_5: user.AccessType,
	}
	q := db.New(database.DB)
	err := q.ShareSecret(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return
	}
}
func GoSliceToPostgresArray(arr []string) string {
	return "{" + strings.Join(arr, ",") + "}"
}

func FetchCredentialByID(ctx *gin.Context, credentialID uuid.UUID) (db.GetCredentialDetailsRow, error) {
	q := db.New(database.DB)
	credential, err := q.GetCredentialDetails(ctx, credentialID)
	if err != nil {
		logger.Errorf(err.Error())
		return credential, err
	}
	return credential, nil

}

func FetchEncryptedData(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) ([]db.GetUserEncryptedDataRow, error) {
	q := db.New(database.DB)
	arg := db.GetUserEncryptedDataParams{
		CredentialID: uuid.NullUUID{UUID: credentialID, Valid: true},
		UserID:       uuid.NullUUID{UUID: userID, Valid: true},
	}
	encryptedData, err := q.GetUserEncryptedData(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	return encryptedData, err
}

func FetchUnEncryptedData(ctx *gin.Context, credentialID uuid.UUID) ([]db.GetCredentialUnencryptedDataRow, error) {
	q := db.New(database.DB)

	encryptedData, err := q.GetCredentialUnencryptedData(ctx, uuid.NullUUID{UUID: credentialID, Valid: true})
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	return encryptedData, err
}
