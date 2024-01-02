package service

import (
	"osvauld/infra/logger"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var accessLevels = map[string]int{
	"unauthorized": 0,
	"read":         1,
	"write":        2,
	"owner":        99,
}

func GetAccessForCredential(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) (string, error) {

	accessRows, err := repository.GetCredentialAccessForUser(ctx, credentialID, userID)
	if err != nil {
		return "", err
	}

	// Find the higest access level for a user
	highestAccess := "unauthorized"
	for _, accessRow := range accessRows {

		if accessLevels[accessRow.AccessType] > accessLevels[highestAccess] {
			highestAccess = accessRow.AccessType
		}
	}

	return highestAccess, nil

}

func HasReadAccessForCredential(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) (bool, error) {
	access, err := GetAccessForCredential(ctx, credentialID, userID)
	logger.Debugf("access level for credential %s, user: %s is %s", credentialID, userID, access)
	if err != nil {
		return false, err
	}
	if accessLevels[access] > 0 {
		return true, nil
	}
	return false, nil
}

func HasOwnerAccessForCredential(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) (bool, error) {

	access, err := GetAccessForCredential(ctx, credentialID, userID)
	logger.Debugf("access level for credential %s, user: %s is %s", credentialID, userID, access)
	if err != nil {
		return false, err
	}
	if accessLevels[access] == 99 {
		return true, nil
	}
	return false, nil

}
