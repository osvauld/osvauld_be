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

var CredentialAccessLevels = map[string]int{
	"unauthorized": 0,
	"read":         1,
	"write":        2,
	"owner":        99,
}

func GetAccessTypeForCredential(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) (string, error) {

	accessRows, err := repository.GetCredentialAccessForUser(ctx, credentialID, userID)
	if err != nil {
		return "", err
	}

	// Find the higest access level for a user
	highestAccess := "unauthorized"
	for _, accessRow := range accessRows {

		if CredentialAccessLevels[accessRow.AccessType] > CredentialAccessLevels[highestAccess] {
			highestAccess = accessRow.AccessType
		}
	}

	return highestAccess, nil

}

func HasReadAccessForCredential(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) (bool, error) {
	access, err := GetAccessTypeForCredential(ctx, credentialID, userID)
	logger.Debugf("access level for credential %s, user: %s is %s", credentialID, userID, access)
	if err != nil {
		return false, err
	}
	if CredentialAccessLevels[access] > 0 {
		return true, nil
	}
	return false, nil
}

func CheckHasReadAccessForCredential(ctx *gin.Context, access string) bool {
	return CredentialAccessLevels[access] > 0

}

func HasOwnerAccessForCredential(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) (bool, error) {

	access, err := GetAccessTypeForCredential(ctx, credentialID, userID)
	logger.Debugf("access level for credential %s, user: %s is %s", credentialID, userID, access)
	if err != nil {
		return false, err
	}
	if CredentialAccessLevels[access] == 99 {
		return true, nil
	}
	return false, nil

}

func HasOwnerAccessForCredentials(ctx *gin.Context, credentialIDs []uuid.UUID, userID uuid.UUID) (bool, error) {

	// TODO: optimize this
	for _, credentialID := range credentialIDs {

		isOwner, err := HasOwnerAccessForCredential(ctx, credentialID, userID)
		if err != nil {
			return false, err
		}
		if !isOwner {
			errMsg := fmt.Sprintf("user %s does not have owner access for credential %s", userID, credentialID)
			return false, &customerrors.UserNotAnOwnerOfCredentialError{Message: errMsg}
		}
	}

	return true, nil

}

func RemoveCredentialAccessForUsers(ctx *gin.Context, payload dto.RemoveCredentialAccessForUsers, caller uuid.UUID) error {

	// Check caller has owner access for credential
	isOwner, err := HasOwnerAccessForCredentials(ctx, []uuid.UUID{payload.CredentialID}, caller)
	if err != nil {
		return err
	}

	if !isOwner {
		errMsg := fmt.Sprintf("user %s does not have owner access for credential %s", caller, payload.CredentialID)
		return &customerrors.UserNotAnOwnerOfCredentialError{Message: errMsg}
	}

	err = repository.RemoveCredentialAccessForUser(ctx, db.RemoveCredentialAccessForUsersParams{
		UserIds:      payload.UserIDs,
		CredentialID: payload.CredentialID,
	})
	if err != nil {
		return err
	}

	return nil

}
