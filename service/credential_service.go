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

func AddCredential(ctx *gin.Context, request dto.AddCredentialRequest, caller uuid.UUID) error {

	// Retrieve access types for the folder
	accessList, err := repository.GetFolderAccess(ctx, request.FolderID)
	if err != nil {
		return err
	}

	payload := dto.AddCredentialDto{
		AddCredentialRequest: request,
	}
	accessTypeMap := make(map[uuid.UUID]string)

	for _, access := range accessList {
		if access.UserID != caller {
			accessTypeMap[access.UserID] = access.AccessType

		} else {
			//caller should be owner
			accessTypeMap[access.UserID] = "owner"
		}
	}

	/* Convert UserEncryptedFields to UserEncryptedFieldsWithAccess
	 * access type is derived from folder ownership
	 */
	var extendedFields []dto.EncryptedFieldWithAccess
	for _, field := range request.UserEncryptedFields {
		accessType, exists := accessTypeMap[field.UserID]
		if !exists {
			continue
		}
		extendedField := dto.EncryptedFieldWithAccess{
			AddCredentialEncryptedField: field,
			AccessType:                  accessType,
		}
		extendedFields = append(extendedFields, extendedField)
	}
	payload.UserEncryptedFieldsWithAccess = extendedFields

	if err != nil {
		return err
	}
	_, err = repository.AddCredential(ctx, payload, caller)
	if err != nil {
		return err
	}
	return nil
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
