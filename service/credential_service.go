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

func ShareCredentialWithUser(ctx *gin.Context, credentialID uuid.UUID, payload dto.UserEncryptedData, caller uuid.UUID) error {

	hasAccess, err := HasOwnerAccessForCredential(ctx, credentialID, caller)
	if err != nil {
		return err
	}

	if !hasAccess {
		return &customerrors.UserNotAuthenticatedError{Message: "user does not have share access to the credential"}
	}

	return repository.ShareCredentialWithUser(ctx, credentialID, payload)
}

func ShareMultipleCredentialsWithMultipleUsers(ctx *gin.Context, payload []dto.ShareCredentialWithUsers, caller uuid.UUID) ([]map[string]interface{}, error) {

	responses := make([]map[string]interface{}, 0)
	for _, credentialData := range payload {

		credentialShareResponse := make(map[string]interface{})
		credentialShareResponse["credentialId"] = credentialData.CredentialID
		credentialShareResponse["users"] = make([]map[string]interface{}, 0)

		for _, userData := range credentialData.UserEncryptedData {

			userSharedResponse := make(map[string]interface{})
			err := ShareCredentialWithUser(ctx, credentialData.CredentialID, userData, caller)

			if err != nil {
				userSharedResponse["userId"] = userData.UserID
				userSharedResponse["status"] = "failed"
				userSharedResponse["message"] = err.Error()
				users := credentialShareResponse["users"].([]map[string]interface{})
				credentialShareResponse["users"] = append(users, userSharedResponse)

			} else {
				userSharedResponse["userId"] = userData.UserID
				userSharedResponse["status"] = "success"
				userSharedResponse["message"] = "shared successfully"
				users := credentialShareResponse["users"].([]map[string]interface{})
				credentialShareResponse["users"] = append(users, userSharedResponse)
			}
		}
		responses = append(responses, credentialShareResponse)
	}

	return responses, nil
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
