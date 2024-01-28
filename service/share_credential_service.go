package service

import (
	db "osvauld/db/sqlc"
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

		userFieldRecords := []db.AddFieldDataParams{}
		credentialAccessRecords := []db.AddCredentialAccessParams{}

		for _, credential := range userData.CredentialData {

			for _, field := range credential.Fields {

				fieldDetails, err := FetchFieldNameAndTypeByFieldIDForUser(ctx, field.ID, caller)
				if err != nil {
					return nil, err
				}
				fieldRecord := db.AddFieldDataParams{
					FieldName:    fieldDetails.FieldName,
					FieldValue:   field.FieldValue,
					FieldType:    fieldDetails.FieldType,
					CredentialID: credential.CredentialID,
					UserID:       userData.UserID,
				}
				userFieldRecords = append(userFieldRecords, fieldRecord)
			}

			credentialAccessRecord := db.AddCredentialAccessParams{
				CredentialID: credential.CredentialID,
				UserID:       userData.UserID,
				AccessType:   userData.AccessType,
			}
			credentialAccessRecords = append(credentialAccessRecords, credentialAccessRecord)
		}

		userShareResponse := ShareCredentialsWithUserResponse{
			UserID: userData.UserID,
		}
		// Share all the credentials for a user in a single transaction
		shareCredentialParams := db.ShareCredentialTransactionParams{
			FieldArgs:            userFieldRecords,
			CredentialAccessArgs: credentialAccessRecords,
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

		groupFieldRecords := []db.AddFieldDataParams{}
		credentialAccessRecords := []db.AddCredentialAccessParams{}
		for _, userData := range groupData.UserData {

			for _, credential := range userData.CredentialData {

				for _, field := range credential.Fields {

					fieldDetails, err := FetchFieldNameAndTypeByFieldIDForUser(ctx, field.ID, caller)
					if err != nil {
						return nil, err
					}

					fieldRecord := db.AddFieldDataParams{
						FieldName:    fieldDetails.FieldName,
						FieldValue:   field.FieldValue,
						FieldType:    fieldDetails.FieldType,
						CredentialID: credential.CredentialID,
						UserID:       userData.UserID,
					}

					groupFieldRecords = append(groupFieldRecords, fieldRecord)
				}

				credentialAccessRecord := db.AddCredentialAccessParams{
					CredentialID: credential.CredentialID,
					UserID:       userData.UserID,
					AccessType:   userData.AccessType,
					GroupID:      uuid.NullUUID{UUID: groupData.GroupID, Valid: true},
				}
				credentialAccessRecords = append(credentialAccessRecords, credentialAccessRecord)

			}
		}

		groupShareResponse := ShareCredentialsWithGroupResponse{}
		groupShareResponse.GroupID = groupData.GroupID

		shareCredentialParams := db.ShareCredentialTransactionParams{
			FieldArgs:            groupFieldRecords,
			CredentialAccessArgs: credentialAccessRecords,
		}

		err := repository.ShareCredentials(ctx, shareCredentialParams)
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

		userFieldRecords := []db.AddFieldDataParams{}
		credentialAccessRecords := []db.AddCredentialAccessParams{}
		folderAccessRecords := []db.AddFolderAccessParams{}

		for _, credential := range userData.CredentialData {

			for _, field := range credential.Fields {

				fieldDetails, err := FetchFieldNameAndTypeByFieldIDForUser(ctx, field.ID, caller)
				if err != nil {
					return nil, err
				}
				fieldRecord := db.AddFieldDataParams{
					FieldName:    fieldDetails.FieldName,
					FieldValue:   field.FieldValue,
					FieldType:    fieldDetails.FieldType,
					CredentialID: credential.CredentialID,
					UserID:       userData.UserID,
				}
				userFieldRecords = append(userFieldRecords, fieldRecord)
			}

			credentialAccessRecord := db.AddCredentialAccessParams{
				CredentialID: credential.CredentialID,
				UserID:       userData.UserID,
				AccessType:   userData.AccessType,
			}
			credentialAccessRecords = append(credentialAccessRecords, credentialAccessRecord)
		}

		folderAccessRecord := db.AddFolderAccessParams{
			UserID:     userData.UserID,
			AccessType: userData.AccessType,
			FolderID:   payload.FolderID,
		}
		folderAccessRecords = append(folderAccessRecords, folderAccessRecord)

		userShareResponse := ShareFolderWithUserResponse{}
		userShareResponse.UserID = userData.UserID

		// Share all the credentials for a user in a single transaction

		shareCredentialTransactionParams := db.ShareCredentialTransactionParams{
			FieldArgs:            userFieldRecords,
			CredentialAccessArgs: credentialAccessRecords,
			FolderAccessArgs:     folderAccessRecords,
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

		groupFieldRecords := []db.AddFieldDataParams{}
		credentialAccessRecords := []db.AddCredentialAccessParams{}
		folderAccessRecords := []db.AddFolderAccessParams{}

		for _, userData := range groupData.UserData {

			for _, credential := range userData.CredentialData {

				for _, field := range credential.Fields {

					fieldDetails, err := FetchFieldNameAndTypeByFieldIDForUser(ctx, field.ID, caller)
					if err != nil {
						return nil, err
					}
					fieldRecord := db.AddFieldDataParams{
						FieldName:    fieldDetails.FieldName,
						FieldValue:   field.FieldValue,
						FieldType:    fieldDetails.FieldType,
						CredentialID: credential.CredentialID,
						UserID:       userData.UserID,
					}
					groupFieldRecords = append(groupFieldRecords, fieldRecord)

				}

				credentialAccessRecord := db.AddCredentialAccessParams{
					CredentialID: credential.CredentialID,
					UserID:       userData.UserID,
					AccessType:   userData.AccessType,
					GroupID:      uuid.NullUUID{UUID: groupData.GroupID, Valid: true},
				}
				credentialAccessRecords = append(credentialAccessRecords, credentialAccessRecord)

			}

			folderAccessRecord := db.AddFolderAccessParams{
				UserID:     userData.UserID,
				AccessType: userData.AccessType,
				FolderID:   payload.FolderID,
				GroupID:    uuid.NullUUID{UUID: groupData.GroupID, Valid: true},
			}
			folderAccessRecords = append(folderAccessRecords, folderAccessRecord)
		}

		groupShareResponse := ShareFolderWithGroupResponse{}
		groupShareResponse.GroupID = groupData.GroupID

		shareCredentialTransactionParams := db.ShareCredentialTransactionParams{
			FieldArgs:            groupFieldRecords,
			CredentialAccessArgs: credentialAccessRecords,
			FolderAccessArgs:     folderAccessRecords,
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
