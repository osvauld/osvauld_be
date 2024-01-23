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

func AddCredential(ctx *gin.Context, request dto.AddCredentialRequest, caller uuid.UUID) (uuid.UUID, error) {

	isOwner, err := CheckFolderOwner(ctx, request.FolderID, caller)
	if err != nil {
		return uuid.UUID{}, err
	}
	if !isOwner {
		return uuid.UUID{}, &customerrors.UserNotAuthenticatedError{Message: "user does not have owner access to the folder"}
	}

	// Retrieve access types for the folder
	accessList, err := repository.GetFolderAccess(ctx, request.FolderID)
	if err != nil {
		return uuid.UUID{}, err
	}

	accessTypeMap := make(map[uuid.UUID]string)

	for _, access := range accessList {
		if access.UserID != caller {
			accessTypeMap[access.UserID] = access.AccessType

			// TODO: this check is redundant since only owner can add credentials
		} else {
			//caller should be owner
			accessTypeMap[access.UserID] = "owner"
		}
	}

	/* Convert UserFields to UserFieldsWithAccessType
	 * access type is derived from the users forlder access
	 */
	var UserFieldsWithAccessTypeSlice []dto.UserFieldsWithAccessType
	for _, userFields := range request.UserFields {
		accessType, exists := accessTypeMap[userFields.UserID]
		if !exists {
			// TODO: send appropriate error
			continue
		}
		userFieldsWithAccessType := dto.UserFieldsWithAccessType{
			UserID:     userFields.UserID,
			Fields:     userFields.Fields,
			AccessType: accessType,
		}

		UserFieldsWithAccessTypeSlice = append(UserFieldsWithAccessTypeSlice, userFieldsWithAccessType)
	}

	payload := dto.AddCredentialDto{
		Name:                     request.Name,
		Description:              request.Description,
		FolderID:                 request.FolderID,
		CredentialType:           request.CredentialType,
		UserFieldsWithAccessType: UserFieldsWithAccessTypeSlice,
	}

	credentialID, err := repository.AddCredential(ctx, payload, caller)
	if err != nil {
		return uuid.UUID{}, err
	}
	return credentialID, nil
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

func GetCredentialsByFolder(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) ([]dto.CredentialForUser, error) {
	credentials, err := repository.GetCredentialsByFolder(ctx, folderID, userID)
	if err != nil {
		return nil, err
	}
	return credentials, nil
}

func GetCredentialsFieldsByFolderID(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) ([]db.GetCredentialsFieldsByIdsRow, error) {
	credentialsIds, err := repository.GetCredentialIdsByFolderAndUserId(ctx, folderID, userID)
	if err != nil {
		return nil, err
	}
	credentials, _ := repository.GetCredentialsFieldsByIds(ctx, credentialsIds, userID)
	return credentials, nil
}

func GetCredentialsFieldsByIds(ctx *gin.Context, credentialIds []uuid.UUID, userID uuid.UUID) ([]db.GetCredentialsFieldsByIdsRow, error) {
	credentials, err := repository.GetCredentialsFieldsByIds(ctx, credentialIds, userID)
	if err != nil {
		return nil, err
	}
	return credentials, nil
}

func GetCredentialsByUrl(ctx *gin.Context, url string, userID uuid.UUID) ([]db.GetCredentialDetailsByIdsRow, error) {
	credentials, err := repository.GetCredentialsByUrl(ctx, url, userID)
	if err != nil {
		return nil, err
	}
	return credentials, nil
}

func GetAllUrlsForUser(ctx *gin.Context, userID uuid.UUID) ([]string, error) {
	urls, err := repository.GetAllUrlsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func GetSensitiveFieldsCredentialByID(ctx *gin.Context, credentialID uuid.UUID, caller uuid.UUID) ([]db.GetSensitiveFieldsRow, error) {
	// Check if caller has access
	hasAccess, err := HasReadAccessForCredential(ctx, credentialID, caller)
	var sensitiveFields []db.GetSensitiveFieldsRow
	if err != nil {
		return sensitiveFields, err
	}

	if !hasAccess {
		logger.Errorf("user %s does not have access to the credential %s", caller, credentialID)
		return sensitiveFields, &customerrors.UserNotAuthenticatedError{Message: "user does not have access to the credential"}
	}

	sensitiveFields, err = repository.GetSensitiveFieldsById(ctx, credentialID, caller)

	return sensitiveFields, err
}
