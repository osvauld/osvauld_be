package service

import (
	"fmt"
	"osvauld/customerrors"
	dto "osvauld/dtos"
	"osvauld/repository"
	"osvauld/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// This is the service layer function used when multiple credentials are shared with multiple users
// We will try to insert all the credentials for a single user in a single transaction
// so that we can rollback all the credentials if one of them fails to be shared
// The response contains success or failure for each user
func ShareMultipleCredentialsWithMultipleUsers(ctx *gin.Context, payload []dto.MulitpleCredentialsForUserPayload, caller uuid.UUID) ([]map[string]interface{}, error) {

	uniqueCredentialIDs := []uuid.UUID{}
	userCredentials := map[uuid.UUID][]dto.CredentialEncryptedFieldsForUserDto{}

	// the following loop is for grouping the credentials shared for a sigle user
	// so that we can share all the credentials for a single user in a single transaction
	for _, userData := range payload {

		for _, credential := range userData.CredentialData {
			// if the credential id does not exist in the uniqueCredentialIDs slice add it
			if !utils.Contains(uniqueCredentialIDs, credential.CredentialID) {
				uniqueCredentialIDs = append(uniqueCredentialIDs, credential.CredentialID)
			}

			credentialDataParams := dto.CredentialEncryptedFieldsForUserDto{
				CredentialID:    credential.CredentialID,
				UserID:          userData.UserID,
				AccessType:      userData.AccessType,
				EncryptedFields: credential.EncryptedData,
			}

			userCredentials[userData.UserID] = append(userCredentials[userData.UserID], credentialDataParams)
		}
	}

	// check the service calls has access to the credentials being shared
	for _, credentialID := range uniqueCredentialIDs {
		hasAccess, err := HasOwnerAccessForCredential(ctx, credentialID, caller)
		if err != nil {
			return nil, err
		}

		if !hasAccess {
			return nil, &customerrors.UserNotAuthenticatedError{Message: fmt.Sprintf("user does not have share access to the credential: %s", credentialID)}

		}
	}

	// TODO: Check the number of encrypted fields and field names for each encrypted value match with the original


	responses := make([]map[string]interface{}, 0)
	for userID, credentials := range userCredentials {

		userShareResponse := make(map[string]interface{})
		userShareResponse["userId"] = userID

		// Share all the credentials for a user in a single transaction
		err := repository.ShareMultipleCredentialsWithMultipleUsers(ctx, credentials)
		if err != nil {
			userShareResponse["status"] = "failed"
			userShareResponse["message"] = err.Error()
		} else {
			userShareResponse["status"] = "success"
			userShareResponse["message"] = "shared successfully"
		}
		responses = append(responses, userShareResponse)
	}

	return responses, nil
}

// This is the service layer function used when multiple credentials are shared with multiple groups
// We will try to insert all the credentials for a single group in a single transaction
// so that we can rollback all the credentials if one of them fails to be shared
// The response contains success or failure for each user
func ShareMultipleCredentialsWithMulitpleGroups(ctx *gin.Context, payload []dto.MulitpleCredentialsForGroupPayload, caller uuid.UUID) ([]map[string]interface{}, error) {

	uniqueCredentialIDs := []uuid.UUID{}
	groupCredentials := map[uuid.UUID][]dto.CredentialEncryptedFieldsForGroupDto{}

	// the following loop is for grouping the credentials shared for a sigle user
	// so that we can share all the credentials for a single user in a single transaction
	for _, groupData := range payload {

		for _, credential := range groupData.CredentialData {
			// if the credential id does not exist in the uniqueCredentialIDs slice add it
			if !utils.Contains(uniqueCredentialIDs, credential.CredentialID) {
				uniqueCredentialIDs = append(uniqueCredentialIDs, credential.CredentialID)
			}

			credentialDataParams := dto.CredentialEncryptedFieldsForGroupDto{
				CredentialID:        credential.CredentialID,
				GroupID:             groupData.GroupID,
				AccessType:          groupData.AccessType,
				UserEncryptedFields: credential.EncryptedData,
			}

			groupCredentials[groupData.GroupID] = append(groupCredentials[groupData.GroupID], credentialDataParams)
		}
	}

	for _, credentialID := range uniqueCredentialIDs {
		hasAccess, err := HasOwnerAccessForCredential(ctx, credentialID, caller)
		if err != nil {
			return nil, err
		}

		if !hasAccess {
			return nil, &customerrors.UserNotAuthenticatedError{Message: fmt.Sprintf("user does not have share access to the credential: %s", credentialID)}

		}
	}

	responses := make([]map[string]interface{}, 0)
	for groupID, credentials := range groupCredentials {

		groupShareResponse := make(map[string]interface{})
		groupShareResponse["groupId"] = groupID

		// Share all the credentials for a group in a single transaction
		err := repository.ShareMultipleCredentialsWithMultipleGroups(ctx, credentials)
		if err != nil {
			groupShareResponse["status"] = "failed"
			groupShareResponse["message"] = err.Error()
		} else {
			groupShareResponse["status"] = "success"
			groupShareResponse["message"] = "shared successfully"
		}

		responses = append(responses, groupShareResponse)
	}

	return responses, nil
}
