package service

import (
	"osvauld/customerrors"
	db "osvauld/db/sqlc"
	"osvauld/infra/logger"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetSensitiveFields(ctx *gin.Context, credentialID uuid.UUID, caller uuid.UUID) ([]db.GetSensitiveFieldsRow, error) {
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

	sensitiveFields, err = repository.GetSensitiveFields(ctx, db.GetSensitiveFieldsParams{
		CredentialID: credentialID,
		UserID:       caller,
	})

	return sensitiveFields, err
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

func DeleteAccessRemovedFields(ctx *gin.Context) error {
	return repository.DeleteAccessRemovedFields(ctx)
}
