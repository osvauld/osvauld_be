package service

import (
	"fmt"
	"osvauld/customerrors"
	dto "osvauld/dtos"
	"osvauld/infra/logger"
	"osvauld/repository"
	"osvauld/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// This is the service layer function used when multiple credentials are shared with multiple users
// We will try to insert all the credentials for a single user in a single transaction
// so that we can rollback all the credentials if one of them fails to be shared
// The response contains success or failure for each user
func ShareCredentialsWithUsers(ctx *gin.Context, payload []dto.ShareCredentialsForUserPayload, caller uuid.UUID) ([]map[string]interface{}, error) {

	uniqueCredentialIDs := []uuid.UUID{}

	// get all the unique credential ids
	for _, userData := range payload {
		for _, credential := range userData.CredentialData {
			// if the credential id does not exist in the uniqueCredentialIDs slice add it
			if !utils.Contains(uniqueCredentialIDs, credential.CredentialID) {
				uniqueCredentialIDs = append(uniqueCredentialIDs, credential.CredentialID)
			}
		}
	}

	// check the service calls has access to the credentials being shared
	// TODO: Change this to make a single database call for all credentials
	for _, credentialID := range uniqueCredentialIDs {
		hasAccess, err := HasOwnerAccessForCredential(ctx, credentialID, caller)
		if err != nil {
			return nil, err
		}

		if !hasAccess {
			return nil, &customerrors.UserNotAuthenticatedError{Message: fmt.Sprintf("user does not have share access to the credential: %s", credentialID)}

		}
	}

	userCredentials := map[uuid.UUID][]dto.CredentialFieldsForUserDto{}

	// the following loop is for grouping the credentials shared for a sigle user
	// so that we can share all the credentials for a single user in a single transaction
	for _, userData := range payload {
		for _, credential := range userData.CredentialData {

			credentialFields := []dto.Field{}
			for _, field := range credential.Fields {

				fieldDetails, err := FetchFieldNameAndTypeByFieldIDForUser(ctx, field.ID, caller)
				if err != nil {
					logger.Errorf("error fetching field details: %s", err.Error())
					return nil, err
				}
				fieldDetails.FieldValue = field.FieldValue
				credentialFields = append(credentialFields, fieldDetails)
			}

			credentialDataParams := dto.CredentialFieldsForUserDto{
				CredentialID: credential.CredentialID,
				UserID:       userData.UserID,
				AccessType:   userData.AccessType,
				Fields:       credentialFields,
			}

			userCredentials[userData.UserID] = append(userCredentials[userData.UserID], credentialDataParams)
		}
	}

	// TODO: Check the number of encrypted fields and field names for each encrypted value match with the original

	responses := make([]map[string]interface{}, 0)
	for userID, credentials := range userCredentials {

		userShareResponse := make(map[string]interface{})
		userShareResponse["userId"] = userID

		// Share all the credentials for a user in a single transaction
		err := repository.ShareCredentialsWithUsers(ctx, credentials)
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

type ShareCredentialsWithGroupResponse struct {
	Status  string    `json:"status"`
	Message string    `json:"message,omitempty"`
	GroupID uuid.UUID `json:"groupId"`
}

func ShareCredentialsWithGroups(ctx *gin.Context, payload []dto.CredentialsForGroupsPayload, caller uuid.UUID) []ShareCredentialsWithGroupResponse {
	// TODO: check caller access to folder
	var responses []ShareCredentialsWithGroupResponse
	for _, groupData := range payload {
		err := repository.ShareCredentialWithGroup(ctx, groupData.GroupID, groupData.AccessType, groupData.EncryptedUserData)
		response := ShareCredentialsWithGroupResponse{
			GroupID: groupData.GroupID,
		}
		if err != nil {
			response.Status = "failed"
			response.Message = err.Error()
		} else {
			response.Status = "success"
			response.Message = "shared successfully"
		}
		responses = append(responses, response)
	}
	return responses
}

func ShareFolderWithUsers(ctx *gin.Context, folderId uuid.UUID, encryptedPayload []dto.CredentialsForUsersPayload) error {
	//TODO: modify the payload to add only to access list or add to encrypted table and access list
	// TODO: make this idempotent
	for _, credentials := range encryptedPayload {
		err := repository.ShareFolderWithUsers(ctx, folderId, credentials)
		if err != nil {
			logger.Debugf("error sharing folder with user: %s", err.Error())
			return err
		}
	}
	return nil

}

func ShareFolderWithGroups(ctx *gin.Context, folderId uuid.UUID, encryptedPayload []dto.CredentialsForGroupsPayload) error {

	//TODO: modify the payload to add only to access list or add to encrypted table and access list
	// TODO: make this idempotent
	for _, groupData := range encryptedPayload {
		// TODO: response object with success true or false
		err := repository.ShareFolderWithGroup(ctx, folderId, groupData)
		if err != nil {
			return err
		}
	}
	return nil

}
