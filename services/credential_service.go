package service

import (
	"errors"
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/logger"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateCredential(ctx *gin.Context, data dto.AddCredentailRequest, userID uuid.UUID) error {
	id, err := repository.SaveCredential(ctx, data.Name, data.Description, data.FolderID, userID)
	for _, encryptedData := range data.EncryptedFields {

		repository.SaveEncryptedData(ctx, encryptedData, id, userID)
	}
	for _, unencryptedData := range data.UnencryptedFields {
		repository.SaveUnEncryptedData(ctx, unencryptedData, id)
	}
	repository.AddToAccessList(ctx, id, "owner", userID)

	return err
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
		return dto.CredentialDetails{}, errors.New("User does not have access to the credential")
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
