package service

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUniqueCredentialIDs(credentials []dto.ShareCredentialPayload) []uuid.UUID {

	// find all unique credential ids
	credentialIDMap := make(map[uuid.UUID]bool)
	for _, credential := range credentials {
		credentialIDMap[credential.CredentialID] = true
	}

	credentialIDs := []uuid.UUID{}
	for credentialID := range credentialIDMap {
		credentialIDs = append(credentialIDs, credentialID)
	}
	return credentialIDs
}

func CreateFieldDataRecords(ctx *gin.Context, credentials []dto.ShareCredentialPayload, userID uuid.UUID) ([]db.AddFieldDataParams, error) {

	credentialIDs := GetUniqueCredentialIDs(credentials)

	fieldData, err := repository.GetFieldDataForCredentialIDsForUser(ctx, db.GetFieldDataByCredentialIDsForUserParams{
		UserID:      userID,
		Credentials: credentialIDs,
	})
	if err != nil {
		return nil, err
	}

	fieldMap := make(map[uuid.UUID]db.GetFieldDataByCredentialIDsForUserRow)
	for _, field := range fieldData {
		fieldMap[field.ID] = field
	}

	userFieldRecords := []db.AddFieldDataParams{}
	for _, credential := range credentials {

		for _, field := range credential.Fields {

			userFieldRecord := db.AddFieldDataParams{
				FieldName:    fieldMap[field.ID].FieldName,
				FieldType:    fieldMap[field.ID].FieldType,
				FieldValue:   field.FieldValue,
				UserID:       userID,
				CredentialID: credential.CredentialID,
			}

			userFieldRecords = append(userFieldRecords, userFieldRecord)
		}
	}
	return userFieldRecords, nil

}

type CreateCredentialAccessRecordParams struct {
	CredentialIDs []uuid.UUID
	UserID        uuid.UUID
	AccessType    string
	GroupID       uuid.NullUUID
}

func CreateCredentialAccessRecords(ctx *gin.Context, params CreateCredentialAccessRecordParams) ([]db.AddCredentialAccessParams, error) {

	credentialAccessRecords := []db.AddCredentialAccessParams{}
	for _, credentialID := range params.CredentialIDs {

		credentialAccessRecord := db.AddCredentialAccessParams{
			CredentialID: credentialID,
			UserID:       params.UserID,
			AccessType:   params.AccessType,
			GroupID:      params.GroupID,
		}
		credentialAccessRecords = append(credentialAccessRecords, credentialAccessRecord)
	}
	return credentialAccessRecords, nil
}

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

	responses := []ShareCredentialsWithUserResponse{}

	// we share all the credentials for a single user in a single transaction
	for _, userData := range payload {

		userFieldRecords, err := CreateFieldDataRecords(ctx, userData.CredentialData, userData.UserID)
		if err != nil {
			return nil, err
		}

		credentialIDs := GetUniqueCredentialIDs(userData.CredentialData)
		credentialAccessRecords, err := CreateCredentialAccessRecords(ctx, CreateCredentialAccessRecordParams{
			CredentialIDs: credentialIDs,
			UserID:        userData.UserID,
			AccessType:    userData.AccessType,
		})

		if err != nil {
			return nil, err
		}

		userShareResponse := ShareCredentialsWithUserResponse{
			UserID: userData.UserID,
		}
		// Share all the credentials for a user in a single transaction
		shareCredentialParams := db.ShareCredentialTransactionParams{
			FieldArgs:            userFieldRecords,
			CredentialAccessArgs: credentialAccessRecords,
		}

		err = repository.ShareCredentials(ctx, shareCredentialParams)
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

	// combine credentials for a single group

	var responses []ShareCredentialsWithGroupResponse
	for _, groupData := range payload {

		groupFieldRecords := []db.AddFieldDataParams{}
		credentialAccessRecords := []db.AddCredentialAccessParams{}

		for _, userData := range groupData.UserData {

			userFieldRecords, err := CreateFieldDataRecords(ctx, userData.CredentialData, userData.UserID)
			if err != nil {
				return nil, err
			}

			groupFieldRecords = append(groupFieldRecords, userFieldRecords...)

			credentialIDs := GetUniqueCredentialIDs(userData.CredentialData)
			credentialAccessRecordsForUser, err := CreateCredentialAccessRecords(ctx, CreateCredentialAccessRecordParams{
				CredentialIDs: credentialIDs,
				UserID:        userData.UserID,
				AccessType:    groupData.AccessType,
				GroupID:       uuid.NullUUID{UUID: groupData.GroupID, Valid: true},
			})
			if err != nil {
				return nil, err
			}

			credentialAccessRecords = append(credentialAccessRecords, credentialAccessRecordsForUser...)

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

		userFieldRecords, err := CreateFieldDataRecords(ctx, userData.CredentialData, userData.UserID)
		if err != nil {
			return nil, err
		}

		credentialIDs := GetUniqueCredentialIDs(userData.CredentialData)
		credentialAccessRecords, err := CreateCredentialAccessRecords(ctx, CreateCredentialAccessRecordParams{
			CredentialIDs: credentialIDs,
			UserID:        userData.UserID,
			AccessType:    userData.AccessType,
		})
		if err != nil {
			return nil, err
		}

		folderAccessRecords := []db.AddFolderAccessParams{
			{
				UserID:     userData.UserID,
				AccessType: userData.AccessType,
				FolderID:   payload.FolderID,
			},
		}

		userShareResponse := ShareFolderWithUserResponse{}
		userShareResponse.UserID = userData.UserID

		// Share all the credentials for a user in a single transaction

		shareCredentialTransactionParams := db.ShareCredentialTransactionParams{
			FieldArgs:            userFieldRecords,
			CredentialAccessArgs: credentialAccessRecords,
			FolderAccessArgs:     folderAccessRecords,
		}

		err = repository.ShareCredentials(ctx, shareCredentialTransactionParams)
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

			userFieldRecords, err := CreateFieldDataRecords(ctx, userData.CredentialData, userData.UserID)
			if err != nil {
				return nil, err
			}

			groupFieldRecords = append(groupFieldRecords, userFieldRecords...)

			credentialIDs := GetUniqueCredentialIDs(userData.CredentialData)
			credentialAccessRecordsForUser, err := CreateCredentialAccessRecords(ctx, CreateCredentialAccessRecordParams{
				CredentialIDs: credentialIDs,
				UserID:        userData.UserID,
				AccessType:    groupData.AccessType,
				GroupID:       uuid.NullUUID{UUID: groupData.GroupID, Valid: true},
			})
			if err != nil {
				return nil, err
			}

			credentialAccessRecords = append(credentialAccessRecords, credentialAccessRecordsForUser...)

			folderAccessRecord := db.AddFolderAccessParams{
				UserID:     userData.UserID,
				AccessType: groupData.AccessType,
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
