package service

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/logger"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateFieldDataRecords(ctx *gin.Context, credential dto.ShareCredentialPayload, userID uuid.UUID, caller uuid.UUID) ([]db.AddFieldDataParams, error) {

	fieldData, err := repository.GetAllFieldsForCredentialIDsForUser(ctx, db.GetAllFieldsByCredentialIDsForUserParams{
		UserID:      caller,
		Credentials: []uuid.UUID{credential.CredentialID},
	})
	if err != nil {
		return nil, err
	}

	fieldMap := make(map[uuid.UUID]db.GetAllFieldsByCredentialIDsForUserRow)
	for _, field := range fieldData {
		fieldMap[field.ID] = field
	}

	userFieldRecords := []db.AddFieldDataParams{}

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

	return userFieldRecords, nil

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

		userFieldRecords := []db.AddFieldDataParams{}
		credentialAccessRecords := []db.AddCredentialAccessParams{}
		userShareResponse := ShareCredentialsWithUserResponse{
			UserID: userData.UserID,
		}

		for _, credential := range userData.CredentialData {

			// Check credential already shared for user
			exists, err := repository.CheckCredentialAccessEntryExists(ctx, db.CheckCredentialAccessEntryExistsParams{
				UserID:       userData.UserID,
				CredentialID: credential.CredentialID,
			})
			if err != nil {
				return nil, err
			}
			if exists {
				logger.Infof("Credential %s already shared for user %s", credential.CredentialID, userData.UserID)
				continue
			} else {

				credentialAccessRecord := db.AddCredentialAccessParams{
					CredentialID: credential.CredentialID,
					UserID:       userData.UserID,
					AccessType:   userData.AccessType,
				}
				credentialAccessRecords = append(credentialAccessRecords, credentialAccessRecord)

			}

			credentialAccessRecord := db.AddCredentialAccessParams{
				CredentialID: credential.CredentialID,
				UserID:       userData.UserID,
				AccessType:   userData.AccessType,
			}
			credentialAccessRecords = append(credentialAccessRecords, credentialAccessRecord)

			///////////////////////////////////////////////////////////////////////////////////////////
			// Check Fields already shared for user

			fieldDataExists, err := repository.CheckFieldEntryExists(ctx, db.CheckFieldEntryExistsParams{
				UserID:       userData.UserID,
				CredentialID: credential.CredentialID,
			})
			if err != nil {
				return nil, err
			}

			if !fieldDataExists {
				userFields, err := CreateFieldDataRecords(ctx, credential, userData.UserID, caller)
				if err != nil {
					return nil, err
				}

				userFieldRecords = append(userFieldRecords, userFields...)
			} else {
				logger.Infof("Field data already exists for credential %s and user %s", credential.CredentialID, userData.UserID)
			}

			//////////////////////////////////////////////////////////////////////////////////////////////
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
	// combine credentials for a single group

	var responses []ShareCredentialsWithGroupResponse
	for _, groupData := range payload {

		groupFieldRecords := []db.AddFieldDataParams{}
		credentialAccessRecords := []db.AddCredentialAccessParams{}

		for _, userData := range groupData.UserData {

			for _, credential := range userData.CredentialData {

				// Check credential already shared for user
				exists, err := repository.CheckCredentialAccessEntryExists(ctx, db.CheckCredentialAccessEntryExistsParams{
					UserID:       userData.UserID,
					CredentialID: credential.CredentialID,
					GroupID:      uuid.NullUUID{UUID: groupData.GroupID, Valid: true},
				})
				if err != nil {
					return nil, err
				}
				if exists {
					logger.Infof("Credential %s already shared for user %s", credential.CredentialID, userData.UserID)
					continue
				} else {

					credentialAccessRecord := db.AddCredentialAccessParams{
						CredentialID: credential.CredentialID,
						UserID:       userData.UserID,
						AccessType:   userData.AccessType,
						GroupID:      uuid.NullUUID{UUID: groupData.GroupID, Valid: true},
					}
					credentialAccessRecords = append(credentialAccessRecords, credentialAccessRecord)

				}

				credentialAccessRecord := db.AddCredentialAccessParams{
					CredentialID: credential.CredentialID,
					UserID:       userData.UserID,
					AccessType:   userData.AccessType,
					GroupID:      uuid.NullUUID{UUID: groupData.GroupID, Valid: true},
				}
				credentialAccessRecords = append(credentialAccessRecords, credentialAccessRecord)

				fieldDataExists, err := repository.CheckFieldEntryExists(ctx, db.CheckFieldEntryExistsParams{
					UserID:       userData.UserID,
					CredentialID: credential.CredentialID,
				})
				if err != nil {
					return nil, err
				}

				if !fieldDataExists {
					userFieldRecords, err := CreateFieldDataRecords(ctx, credential, userData.UserID, caller)
					if err != nil {
						return nil, err
					}

					groupFieldRecords = append(groupFieldRecords, userFieldRecords...)
				} else {
					logger.Infof("Field data already exists for credential %s and user %s", credential.CredentialID, userData.UserID)
				}

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

		credentialAccessRecords := []db.AddCredentialAccessParams{}
		userFieldRecords := []db.AddFieldDataParams{}
		folderAccessRecords := []db.AddFolderAccessParams{}

		userShareResponse := ShareFolderWithUserResponse{}
		userShareResponse.UserID = userData.UserID

		for _, credential := range userData.CredentialData {

			fieldDataExists, err := repository.CheckFieldEntryExists(ctx, db.CheckFieldEntryExistsParams{
				UserID:       userData.UserID,
				CredentialID: credential.CredentialID,
			})
			if err != nil {
				return nil, err
			}

			if !fieldDataExists {
				userFields, err := CreateFieldDataRecords(ctx, credential, userData.UserID, caller)
				if err != nil {
					return nil, err
				}
				userFieldRecords = append(userFieldRecords, userFields...)
			} else {
				logger.Infof("Field data already exists for credential %s and user %s", credential.CredentialID, userData.UserID)
			}

			// Check credential already shared for user
			exists, err := repository.CheckCredentialAccessEntryExists(ctx, db.CheckCredentialAccessEntryExistsParams{
				UserID:       userData.UserID,
				CredentialID: credential.CredentialID,
				FolderID:     uuid.NullUUID{UUID: payload.FolderID, Valid: true},
			})
			if err != nil {
				return nil, err
			}
			if exists {
				logger.Infof("Credential %s already shared for user %s", credential.CredentialID, userData.UserID)
			} else {

				credentialAccessRecord := db.AddCredentialAccessParams{
					CredentialID: credential.CredentialID,
					UserID:       userData.UserID,
					AccessType:   userData.AccessType,
					FolderID:     uuid.NullUUID{UUID: payload.FolderID, Valid: true},
				}
				credentialAccessRecords = append(credentialAccessRecords, credentialAccessRecord)
			}

		}

		exists, err := repository.CheckFolderAccessEntryExists(ctx, db.CheckFolderAccessEntryExistsParams{
			UserID:   userData.UserID,
			FolderID: payload.FolderID,
		})
		if err != nil {
			return nil, err
		}
		if exists {
			logger.Infof("Folder %s already shared for user %s", payload.FolderID, userData.UserID)
		} else {

			folderAccessRecord := db.AddFolderAccessParams{
				UserID:     userData.UserID,
				AccessType: userData.AccessType,
				FolderID:   payload.FolderID,
			}
			folderAccessRecords = append(folderAccessRecords, folderAccessRecord)
		}

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

			for _, credential := range userData.CredentialData {

				fieldDataExists, err := repository.CheckFieldEntryExists(ctx, db.CheckFieldEntryExistsParams{
					UserID:       userData.UserID,
					CredentialID: credential.CredentialID,
				})
				if err != nil {
					return nil, err
				}

				if !fieldDataExists {
					userFields, err := CreateFieldDataRecords(ctx, credential, userData.UserID, caller)
					if err != nil {
						return nil, err
					}

					groupFieldRecords = append(groupFieldRecords, userFields...)
				} else {
					logger.Infof("Field data already exists for credential %s and user %s", credential.CredentialID, userData.UserID)
				}

				exists, err := repository.CheckCredentialAccessEntryExists(ctx, db.CheckCredentialAccessEntryExistsParams{
					UserID:       userData.UserID,
					CredentialID: credential.CredentialID,
					GroupID:      uuid.NullUUID{UUID: groupData.GroupID, Valid: true},
					FolderID:     uuid.NullUUID{UUID: payload.FolderID, Valid: true},
				})
				if err != nil {
					return nil, err
				}
				if exists {
					logger.Infof("Credential %s already shared for user %s", credential.CredentialID, userData.UserID)
					continue
				} else {

					credentialAccessRecord := db.AddCredentialAccessParams{
						CredentialID: credential.CredentialID,
						UserID:       userData.UserID,
						AccessType:   userData.AccessType,
						GroupID:      uuid.NullUUID{UUID: groupData.GroupID, Valid: true},
						FolderID:     uuid.NullUUID{UUID: payload.FolderID, Valid: true},
					}

					credentialAccessRecords = append(credentialAccessRecords, credentialAccessRecord)
				}
			}

			exists, err := repository.CheckFolderAccessEntryExists(ctx, db.CheckFolderAccessEntryExistsParams{
				UserID:   userData.UserID,
				FolderID: payload.FolderID,
				GroupID:  uuid.NullUUID{UUID: groupData.GroupID, Valid: true},
			})
			if err != nil {
				return nil, err
			}

			if exists {
				logger.Infof("Folder %s already shared for user %s", payload.FolderID, userData.UserID)
			} else {
				folderAccessRecord := db.AddFolderAccessParams{
					UserID:     userData.UserID,
					AccessType: groupData.AccessType,
					FolderID:   payload.FolderID,
					GroupID:    uuid.NullUUID{UUID: groupData.GroupID, Valid: true},
				}
				folderAccessRecords = append(folderAccessRecords, folderAccessRecord)
			}

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
