package service

import (
	"osvauld/customerrors"
	dto "osvauld/dtos"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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

func ShareCredentialWithGroup(ctx *gin.Context, credentialID uuid.UUID, payload dto.CredentialDataForGroup, caller uuid.UUID) error {

	hasAccess, err := HasOwnerAccessForCredential(ctx, credentialID, caller)
	if err != nil {
		return err
	}

	if !hasAccess {
		return &customerrors.UserNotAuthenticatedError{Message: "user does not have share access to the credential"}
	}

	return repository.ShareCredentialWithGroup(ctx, credentialID, payload)
}

func ShareMultipleCredentialsWithMulitpleGroups(ctx *gin.Context, payload []dto.ShareCredentialWithGroups, caller uuid.UUID) ([]map[string]interface{}, error) {

	responses := make([]map[string]interface{}, 0)
	for _, credentialData := range payload {

		credentialShareResponse := make(map[string]interface{})
		credentialShareResponse["credentialId"] = credentialData.CredentialID
		credentialShareResponse["groups"] = make([]map[string]interface{}, 0)

		for _, groupData := range credentialData.GroupData {

			groupSharedResponse := make(map[string]interface{})
			err := ShareCredentialWithGroup(ctx, credentialData.CredentialID, groupData, caller)

			if err != nil {
				groupSharedResponse["groupId"] = groupData.GroupID
				groupSharedResponse["status"] = "failed"
				groupSharedResponse["message"] = err.Error()
				groups := credentialShareResponse["groups"].([]map[string]interface{})
				credentialShareResponse["groups"] = append(groups, groupSharedResponse)

			} else {
				groupSharedResponse["groupId"] = groupData.GroupID
				groupSharedResponse["status"] = "success"
				groupSharedResponse["message"] = "shared successfully"
				groups := credentialShareResponse["groups"].([]map[string]interface{})
				credentialShareResponse["groups"] = append(groups, groupSharedResponse)
			}
		}
		responses = append(responses, credentialShareResponse)
	}

	return responses, nil
}
