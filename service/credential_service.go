package service

import (
	"fmt"
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

	fmt.Println(credential)
	return credential, err
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
		logger.Infof("Sharing credential with id: %s", credential.CredentialID)
		id := credential.CredentialID
		for _, user := range credential.Users {
			repository.ShareCredential(ctx, id, user)
		}
	}
}

// func FetchCredentialByID(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) (dto.CredentialDetails, error) {
// 	if hasAccess, err := repository.CheckAccessForCredential(ctx, credentialID, userID); !hasAccess {
// 		logger.Errorf(err.Error())
// 		logger.Errorf("user does not have access to the credential")
// 		return dto.CredentialDetails{}, err
// 	}
// 	credential, err := repository.FetchCredentialByID(ctx, credentialID)
// 	if err != nil {
// 		logger.Errorf(err.Error())
// 	}
// 	encryptedData, err := repository.FetchEncryptedData(ctx, credentialID, userID)
// 	if err != nil {
// 		logger.Errorf(err.Error())
// 	}
// 	unEncryptedData, err := repository.FetchUnEncryptedData(ctx, credentialID)
// 	if err != nil {
// 		logger.Errorf(err.Error())
// 	}
// 	userList, err := repository.GetUsersByCredential(ctx, credentialID)
// 	if err != nil {
// 		logger.Errorf(err.Error())
// 	}
// 	credentialDetail := dto.CredentialDetails{
// 		Credential:      credential,
// 		EncryptedData:   encryptedData,
// 		UnencryptedData: unEncryptedData,
// 		Users:           userList,
// 	}
// 	return credentialDetail, err

// }

// func extractUniqueUserIDs(encryptedFields []dto.EncryptedFields) ([]uuid.UUID, error) {
// 	userIDMap := make(map[uuid.UUID]bool)
// 	var uniqueUserIDs []uuid.UUID

// 	for _, field := range encryptedFields {
// 		if _, exists := userIDMap[field.UserID]; !exists {
// 			userIDMap[field.UserID] = true
// 			uniqueUserIDs = append(uniqueUserIDs, field.UserID)
// 		}
// 	}

// 	return uniqueUserIDs, nil
// }

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
