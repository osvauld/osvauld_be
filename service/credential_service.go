package service

import (
	"osvauld/customerrors"
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/logger"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddCredential(ctx *gin.Context, data dto.AddCredentailRequest, createdBy uuid.UUID) (uuid.UUID, error) {

	// add access type for users
	for _, user := range data.UserAccessDetails {
		if user.UserID == createdBy {
			user.AccessType = "owner"
		} else {
			user.AccessType = "read"
		}
	}

	return repository.AddCredential(ctx, data, createdBy)
}

func FetchCredentialByID(ctx *gin.Context, credentialID uuid.UUID, caller uuid.UUID) (dto.CredentialDetails, error) {

	// Check if caller has access
	hasAccess, err := HasReadAccessForCredential(ctx, credentialID, caller)
	if err != nil {
		return dto.CredentialDetails{}, err
	}

	if !hasAccess {
		logger.Errorf("user %s does not have access to the credential %s", caller, credentialID)
		return dto.CredentialDetails{}, &customerrors.UserNotAuthenticatedError{Message: "user does not have access to the credential"}
	}

	credential, err := repository.FetchCredentialByID(ctx, credentialID, caller)
	if err != nil {
		return dto.CredentialDetails{}, err
	}

	unencryptedFields, err := repository.FetchUnencryptedFieldsByCredentialID(ctx, credentialID)
	if err != nil {
		return dto.CredentialDetails{}, err
	}

	credential.UnencryptedFields = unencryptedFields

	encryptedFields, err := repository.FetchEncryptedFieldsByCredentialIDByAndUserID(ctx, credentialID, caller)
	if err != nil {
		return dto.CredentialDetails{}, err
	}

	credential.EncryptedFields = encryptedFields

	return credential, err
}

func GetCredentialsByFolder(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) ([]db.FetchCredentialsByUserAndFolderRow, error) {
	credentials, err := repository.GetCredentialsByFolder(ctx, folderID, userID)
	if err != nil {
		return nil, err
	}
	return credentials, nil
}

func GetEncryptedCredentials(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) ([]db.GetEncryptedCredentialsByFolderRow, error) {
	credentials, err := repository.GetEncryptedCredentails(ctx, folderID, userID)
	if err != nil {
		return nil, err
	}
	return credentials, nil
}

func GetEncryptedCredentialsByIds(ctx *gin.Context, credentialIds []uuid.UUID, userID uuid.UUID) ([]db.GetEncryptedDataByCredentialIdsRow, error) {
	credentials, err := repository.GetEncryptedCredentailsByIds(ctx, credentialIds, userID)
	if err != nil {
		return nil, err
	}
	return credentials, nil
}
