package service

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/logger"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateCredential(ctx *gin.Context, data dto.AddCredentailRequest, userID uuid.UUID) (uuid.UUID, error) {
	uniqueUsersIDs, _ := extractUniqueUserIDs(data.EncryptedFields)
	payload := dto.SQLCPayload{
		Name:              data.Name,
		Description:       data.Description,
		UniqueUserIds:     uniqueUsersIDs,
		EncryptedFields:   data.EncryptedFields,
		UnencryptedFields: data.UnencryptedFields,
		FolderID:          data.FolderID,
		CreatedBy:         userID,
	}
	id, err := repository.AddCredential(ctx, payload)
	return id, err
}

func GetCredentialsByFolder(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) ([]db.FetchCredentialsByUserAndFolderRow, error) {
	credentials, err := repository.GetCredentialsByFolder(ctx, folderID, userID)
	if err != nil {
		return nil, err
	}
	return credentials, nil
}

func ShareCredential(ctx *gin.Context, payload dto.ShareCredentialPayload, userID uuid.UUID) {
	for _, credential := range payload.CredentialList {
		id := credential.CredentialID
		for _, user := range credential.Users {
			repository.ShareCredential(ctx, id, user)
		}
	}
}

func FetchCredentialByID(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) (dto.CredentialDetails, error) {
	if hasAccess, err := repository.CheckAccessForCredential(ctx, credentialID, userID); !hasAccess {
		logger.Errorf(err.Error())
		logger.Errorf("user does not have access to the credential")
		return dto.CredentialDetails{}, err
	}
	credential, err := repository.FetchCredentialByID(ctx, credentialID)
	if err != nil {
		logger.Errorf(err.Error())
	}
	encryptedData, err := repository.FetchEncryptedData(ctx, credentialID, userID)
	if err != nil {
		logger.Errorf(err.Error())
	}
	unEncryptedData, err := repository.FetchUnEncryptedData(ctx, credentialID)
	if err != nil {
		logger.Errorf(err.Error())
	}
	credentialDetail := dto.CredentialDetails{
		Credential:      credential,
		EncryptedData:   encryptedData,
		UnencryptedData: unEncryptedData,
	}
	return credentialDetail, err

}

func extractUniqueUserIDs(encryptedFields []dto.EncryptedFields) ([]uuid.UUID, error) {
	userIDMap := make(map[uuid.UUID]bool)
	var uniqueUserIDs []uuid.UUID

	for _, field := range encryptedFields {
		if _, exists := userIDMap[field.UserID]; !exists {
			userIDMap[field.UserID] = true
			uniqueUserIDs = append(uniqueUserIDs, field.UserID)
		}
	}

	return uniqueUserIDs, nil
}
