package service

import (
	dto "osvauld/dtos"
	"osvauld/repository"
	"osvauld/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ShareCredentialsWithUserResponse struct {
	Status  string    `json:"status"`
	Message string    `json:"message,omitempty"`
	UserID  uuid.UUID `json:"userId"`
}

// This is the service layer function used when multiple credentials are shared with multiple users
// We will try to insert all the credentials for a single user in a single transaction
// so that we can rollback all the credentials if one of them fails to be shared
// The response contains success or failure for each user
func ShareCredentialsWithUsers(ctx *gin.Context, payload []dto.ShareCredentialsForUserPayload, caller uuid.UUID) ([]ShareCredentialsWithUserResponse, error) {

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

	_, err := HasOwnerAccessForCredentials(ctx, uniqueCredentialIDs, caller)
	if err != nil {
		return nil, err
	}

	responses := []ShareCredentialsWithUserResponse{}
	// the following loop is for grouping the credentials shared for a single user
	// so that we can share all the credentials for a single user in a single transaction
	for _, userData := range payload {

		userCredentials := []dto.CredentialFieldsForUserDto{}
		for _, credential := range userData.CredentialData {

			credentialFields := []dto.Field{}
			for _, field := range credential.Fields {

				fieldDetails, err := FetchFieldNameAndTypeByFieldIDForUser(ctx, field.ID, caller)
				if err != nil {
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

			userCredentials = append(userCredentials, credentialDataParams)
		}

		userShareResponse := ShareCredentialsWithUserResponse{
			UserID: userData.UserID,
		}
		// Share all the credentials for a user in a single transaction
		shareCredentialParams := dto.ShareCredentialTransactionParams{
			CredentialArgs: userCredentials,
		}
		err := repository.ShareCredentials(ctx, shareCredentialParams)
		if err != nil {
			userShareResponse.Status = "failed"
			userShareResponse.Message = err.Error()
		} else {
			userShareResponse.Status = "success"
			userShareResponse.Message = "shared successfully"
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

func ShareCredentialsWithGroups(ctx *gin.Context, payload []dto.CredentialsForGroupsPayload, caller uuid.UUID) ([]ShareCredentialsWithGroupResponse, error) {

	// TODO: check caller access to credentials, resume logic from ShareCredentialsWithUsers

	// combine credentials for a single group

	var responses []ShareCredentialsWithGroupResponse
	for _, groupData := range payload {

		credentialsForGroup := []dto.CredentialFieldsForUserDto{}
		for _, userData := range groupData.UserData {

			for _, credential := range userData.CredentialData {

				credentialFields := []dto.Field{}
				for _, field := range credential.Fields {

					fieldDetails, err := FetchFieldNameAndTypeByFieldIDForUser(ctx, field.ID, caller)
					if err != nil {
						return nil, err
					}
					fieldDetails.FieldValue = field.FieldValue
					credentialFields = append(credentialFields, fieldDetails)
				}

				credentialFieldsForUser := dto.CredentialFieldsForUserDto{
					CredentialID: credential.CredentialID,
					UserID:       userData.UserID,
					AccessType:   groupData.AccessType,
					Fields:       credentialFields,
					GroupID:      uuid.NullUUID{UUID: groupData.GroupID, Valid: true},
				}

				credentialsForGroup = append(credentialsForGroup, credentialFieldsForUser)
			}
		}

		groupShareResponse := ShareCredentialsWithGroupResponse{}
		groupShareResponse.GroupID = groupData.GroupID

		shareCredentialTransactionParams := dto.ShareCredentialTransactionParams{
			CredentialArgs: credentialsForGroup,
		}

		err := repository.ShareCredentials(ctx, shareCredentialTransactionParams)
		if err != nil {
			groupShareResponse.Status = "failed"
			groupShareResponse.Message = err.Error()
		} else {
			groupShareResponse.Status = "success"
			groupShareResponse.Message = "shared successfully"
		}
		responses = append(responses, groupShareResponse)
	}
	return responses, nil
}

type ShareFolderWithUserResponse struct {
	Status  string    `json:"status"`
	Message string    `json:"message,omitempty"`
	UserID  uuid.UUID `json:"userId"`
}

func ShareFolderWithUsers(ctx *gin.Context, payload dto.ShareFolderWithUsersRequest, caller uuid.UUID) ([]ShareFolderWithUserResponse, error) {
	// TODO: modify the payload to add only to access list or add to encrypted table and access list
	// TODO: make this idempotent

	// the following loop is for grouping the credentials shared for a single user
	// so that we can share all the credentials for a single user in a single transaction
	var responses []ShareFolderWithUserResponse
	for _, userData := range payload.UserData {

		credentialsForUser := []dto.CredentialFieldsForUserDto{}
		for _, credential := range userData.CredentialData {

			credentialFields := []dto.Field{}
			for _, field := range credential.Fields {

				fieldDetails, err := FetchFieldNameAndTypeByFieldIDForUser(ctx, field.ID, caller)
				if err != nil {
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
				FolderID:     uuid.NullUUID{Valid: true, UUID: payload.FolderID},
			}

			credentialsForUser = append(credentialsForUser, credentialDataParams)
		}

		userFolderAccess := dto.UserFolderAccessDto{
			UserID:     userData.UserID,
			AccessType: userData.AccessType,
			FolderID:   payload.FolderID,
		}

		userShareResponse := ShareFolderWithUserResponse{}
		userShareResponse.UserID = userData.UserID

		// Share all the credentials for a user in a single transaction

		shareCredentialTransactionParams := dto.ShareCredentialTransactionParams{
			CredentialArgs: credentialsForUser,
			FolderAccess:   []dto.UserFolderAccessDto{userFolderAccess},
		}

		err := repository.ShareCredentials(ctx, shareCredentialTransactionParams)
		if err != nil {
			userShareResponse.Status = "failed"
			userShareResponse.Message = err.Error()
		} else {
			userShareResponse.Status = "success"
			userShareResponse.Message = "shared successfully"
		}
		responses = append(responses, userShareResponse)
	}

	return responses, nil

}

type ShareFolderWithGroupResponse struct {
	Status  string    `json:"status"`
	Message string    `json:"message,omitempty"`
	GroupID uuid.UUID `json:"groupId"`
}

func ShareFolderWithGroups(ctx *gin.Context, payload dto.ShareFolderWithGroupsRequest, caller uuid.UUID) ([]ShareFolderWithGroupResponse, error) {

	var responses []ShareFolderWithGroupResponse
	for _, groupData := range payload.GroupData {

		credentialsForGroup := []dto.CredentialFieldsForUserDto{}
		userFolderAccesses := []dto.UserFolderAccessDto{}
		for _, userData := range groupData.UserData {

			for _, credential := range userData.CredentialData {

				credentialFields := []dto.Field{}
				for _, field := range credential.Fields {

					fieldDetails, err := FetchFieldNameAndTypeByFieldIDForUser(ctx, field.ID, caller)
					if err != nil {
						return nil, err
					}
					fieldDetails.FieldValue = field.FieldValue
					credentialFields = append(credentialFields, fieldDetails)
				}

				credentialFieldsForUser := dto.CredentialFieldsForUserDto{
					CredentialID: credential.CredentialID,
					UserID:       userData.UserID,
					AccessType:   groupData.AccessType,
					Fields:       credentialFields,
					GroupID:      uuid.NullUUID{UUID: groupData.GroupID, Valid: true},
					FolderID:     uuid.NullUUID{Valid: true, UUID: payload.FolderID},
				}

				userFolderAccess := dto.UserFolderAccessDto{
					UserID:     userData.UserID,
					AccessType: groupData.AccessType,
					FolderID:   payload.FolderID,
					GroupID:    uuid.NullUUID{UUID: groupData.GroupID, Valid: true},
				}

				credentialsForGroup = append(credentialsForGroup, credentialFieldsForUser)
				userFolderAccesses = append(userFolderAccesses, userFolderAccess)
			}
		}

		groupShareResponse := ShareFolderWithGroupResponse{}
		groupShareResponse.GroupID = groupData.GroupID

		shareCredentialTransactionParams := dto.ShareCredentialTransactionParams{
			CredentialArgs: credentialsForGroup,
			FolderAccess:   userFolderAccesses,
		}

		// Share all the credentials for a user in a single transaction
		err := repository.ShareCredentials(ctx, shareCredentialTransactionParams)
		if err != nil {
			groupShareResponse.Status = "failed"
			groupShareResponse.Message = err.Error()
		} else {
			groupShareResponse.Status = "success"
			groupShareResponse.Message = "shared successfully"
		}
		responses = append(responses, groupShareResponse)
	}

	return responses, nil
}
