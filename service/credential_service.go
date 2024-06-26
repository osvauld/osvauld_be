package service

import (
	"database/sql"
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddCredential(ctx *gin.Context, request dto.AddCredentialRequest, caller uuid.UUID) (uuid.UUID, error) {

	if err := VerifyFolderManageAccessForUser(ctx, request.FolderID, caller); err != nil {
		return uuid.UUID{}, err
	}
	// Retrieve access types for the folder
	accessList, err := repository.GetFolderAccess(ctx, request.FolderID)
	if err != nil {
		return uuid.UUID{}, err
	}

	credentialAccessRecords := []db.AddCredentialAccessParams{}
	for _, accessRow := range accessList {

		credentialAccessRecord := db.AddCredentialAccessParams{
			UserID:     accessRow.UserID,
			AccessType: accessRow.AccessType,
			FolderID:   uuid.NullUUID{UUID: request.FolderID, Valid: true},
			GroupID:    accessRow.GroupID,
		}
		credentialAccessRecords = append(credentialAccessRecords, credentialAccessRecord)
	}

	payload := db.AddCredentialTransactionParams{
		Name:                 request.Name,
		Description:          sql.NullString{String: request.Description, Valid: true},
		FolderID:             request.FolderID,
		CredentialType:       request.CredentialType,
		CreatedBy:            caller,
		Fields:               request.Fields,
		CredentialAccessArgs: credentialAccessRecords,
		Domain:               request.Domain,
	}
	credentialID, err := repository.AddCredential(ctx, payload)
	if err != nil {
		return uuid.UUID{}, err
	}
	return credentialID, nil
}

// GetCredentialByID returns the credential details and non sensitive fields for the given credentialID
func GetCredentialDataByID(ctx *gin.Context, credentialID uuid.UUID, caller uuid.UUID) (dto.CredentialForUser, error) {

	if err := VerifyCredentialReadAccessForUser(ctx, credentialID, caller); err != nil {
		return dto.CredentialForUser{}, err
	}

	credential, err := repository.GetCredentialDataByID(ctx, credentialID)
	if err != nil {
		return dto.CredentialForUser{}, err
	}

	fields, err := repository.GetNonSensitiveFieldsForCredentialIDs(ctx, db.GetNonSensitiveFieldsForCredentialIDsParams{
		Credentialids: []uuid.UUID{credentialID},
		UserID:        caller,
	})
	if err != nil {
		return dto.CredentialForUser{}, err
	}

	accessType, err := GetCredentialAccessTypeForUser(ctx, credentialID, caller)
	if err != nil {
		return dto.CredentialForUser{}, err
	}

	fieldDtos := []dto.Field{}
	for _, field := range fields {
		fieldDtos = append(fieldDtos, dto.Field{
			ID:         field.ID,
			FieldName:  field.FieldName,
			FieldValue: field.FieldValue,
			FieldType:  field.FieldType,
		})
	}

	credentialDetails := dto.CredentialForUser{
		CredentialID:   credential.ID,
		Name:           credential.Name,
		Description:    credential.Description.String,
		CredentialType: credential.CredentialType,
		AccessType:     accessType,
		FolderID:       credential.FolderID,
		CreatedAt:      credential.CreatedAt,
		UpdatedAt:      credential.UpdatedAt,
		CreatedBy:      credential.CreatedBy.UUID,
		Fields:         fieldDtos,
	}
	return credentialDetails, err
}

func GetUniqueCredentialsWithHighestAccess(credentials []db.FetchCredentialDetailsForUserByFolderIdRow) []db.FetchCredentialDetailsForUserByFolderIdRow {
	credentialMap := make(map[uuid.UUID]db.FetchCredentialDetailsForUserByFolderIdRow)
	for _, credential := range credentials {
		if _, ok := credentialMap[credential.CredentialID]; ok {

			existingAccessType := credentialMap[credential.CredentialID].AccessType
			newAccessType := credential.AccessType

			if CredentialAccessLevels[newAccessType] > CredentialAccessLevels[existingAccessType] {
				credentialMap[credential.CredentialID] = credential
			}
		} else {
			credentialMap[credential.CredentialID] = credential
		}
	}

	uniqueCredentials := []db.FetchCredentialDetailsForUserByFolderIdRow{}
	for _, credential := range credentialMap {
		uniqueCredentials = append(uniqueCredentials, credential)
	}

	return uniqueCredentials
}

func GetCredentialsByFolder(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) ([]dto.CredentialForUser, error) {

	// Users can have access to only some of the credentials in a folder.
	// So check the access_list table to see which credentials the user has access to
	credentialDetails, err := repository.FetchCredentialDetailsForUserByFolderId(ctx, db.FetchCredentialDetailsForUserByFolderIdParams{
		FolderID: folderID,
		UserID:   userID,
	})
	if err != nil {
		return []dto.CredentialForUser{}, err
	}

	uniqueCredentialDetails := GetUniqueCredentialsWithHighestAccess(credentialDetails)

	credentialIDs := []uuid.UUID{}
	for _, credential := range uniqueCredentialDetails {
		credentialIDs = append(credentialIDs, credential.CredentialID)
	}

	FieldsData, err := repository.GetNonSensitiveFieldsForCredentialIDs(ctx, db.GetNonSensitiveFieldsForCredentialIDsParams{
		Credentialids: credentialIDs,
		UserID:        userID,
	})
	if err != nil {
		return []dto.CredentialForUser{}, err
	}

	credentialFieldGroups := map[uuid.UUID][]dto.Field{}

	for _, field := range FieldsData {
		// if credential.CredentialID does not exist add it to the map and add the field to the array
		credentialFieldGroups[field.CredentialID] = append(credentialFieldGroups[field.CredentialID], dto.Field{
			ID:         field.ID,
			FieldName:  field.FieldName,
			FieldValue: field.FieldValue,
			FieldType:  field.FieldType,
		})
	}

	credentials := []dto.CredentialForUser{}
	for _, credential := range uniqueCredentialDetails {
		credentialForUser := dto.CredentialForUser{}

		credentialForUser.CredentialID = credential.CredentialID
		credentialForUser.Name = credential.Name
		credentialForUser.Description = credential.Description
		credentialForUser.CredentialType = credential.CredentialType
		credentialForUser.AccessType = credential.AccessType
		credentialForUser.FolderID = folderID
		credentialForUser.CreatedAt = credential.CreatedAt
		credentialForUser.UpdatedAt = credential.UpdatedAt
		credentialForUser.CreatedBy = credential.CreatedBy.UUID
		credentialForUser.Fields = credentialFieldGroups[credential.CredentialID]
		if fields, ok := credentialFieldGroups[credential.CredentialID]; ok {
			credentialForUser.Fields = fields
		} else {
			credentialForUser.Fields = []dto.Field{} // Add an empty array if fields are not found
		}

		credentials = append(credentials, credentialForUser)
	}

	return credentials, nil
}

func GetCredentialsByIDs(ctx *gin.Context, credentialIDs []uuid.UUID, userID uuid.UUID) ([]dto.CredentialForUser, error) {

	if err := VerifyReadAccessForCredentials(ctx, credentialIDs, userID); err != nil {
		return nil, err
	}

	FieldsData, err := repository.GetNonSensitiveFieldsForCredentialIDs(ctx, db.GetNonSensitiveFieldsForCredentialIDsParams{
		Credentialids: credentialIDs,
		UserID:        userID,
	})
	if err != nil {
		return nil, err
	}

	credentialFieldGroups := map[uuid.UUID][]dto.Field{}
	for _, field := range FieldsData {
		// if credential.CredentialID does not exist add it to the map and add the field to the array
		credentialFieldGroups[field.CredentialID] = append(credentialFieldGroups[field.CredentialID], dto.Field{
			ID:         field.ID,
			FieldName:  field.FieldName,
			FieldValue: field.FieldValue,
			FieldType:  field.FieldType,
		})
	}

	credentialDetails, err := repository.GetCredentialDetailsByIDs(ctx, credentialIDs)

	credentials := []dto.CredentialForUser{}
	for _, credential := range credentialDetails {
		credentialForUser := dto.CredentialForUser{}

		credentialForUser.CredentialID = credential.ID
		credentialForUser.Name = credential.Name
		credentialForUser.Description = credential.Description.String
		credentialForUser.CredentialType = credential.CredentialType
		credentialForUser.FolderID = credential.FolderID
		credentialForUser.CreatedAt = credential.CreatedAt
		credentialForUser.UpdatedAt = credential.UpdatedAt
		credentialForUser.CreatedBy = credential.CreatedBy.UUID
		credentialForUser.Fields = credentialFieldGroups[credential.ID]

		credentials = append(credentials, credentialForUser)
	}

	if err != nil {
		return nil, err
	}

	return credentials, nil
}

func GetAllUrlsForUser(ctx *gin.Context, userID uuid.UUID) ([]db.GetAllUrlsForUserRow, error) {
	urls, err := repository.GetAllUrlsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return urls, nil
}

func EditCredential(ctx *gin.Context, credentialID uuid.UUID, request dto.EditCredentialRequest, caller uuid.UUID) error {

	if err := VerifyCredentialManageAccessForUser(ctx, credentialID, caller); err != nil {
		return err
	}

	err := repository.EditCredential(ctx, request, caller)
	if err != nil {
		return err
	}

	return nil

}

func EditCredentialDetails(ctx *gin.Context, credentialID uuid.UUID, request dto.EditCredentialDetailsRequest, caller uuid.UUID) error {

	if err := VerifyCredentialManageAccessForUser(ctx, credentialID, caller); err != nil {
		return err
	}

	err := repository.EditCredentialDetails(ctx, db.EditCredentialDetailsParams{
		ID:             credentialID,
		Name:           request.Name,
		Description:    sql.NullString{String: request.Description, Valid: true},
		CredentialType: request.CredentialType,
		UpdatedBy:      uuid.NullUUID{UUID: caller, Valid: true},
	})

	if err != nil {
		return err
	}

	return nil

}

func GetSearchData(ctx *gin.Context, userID uuid.UUID) ([]db.GetCredentialsForSearchByUserIDRow, error) {
	credentials, err := repository.GetSearchData(ctx, userID)
	if err != nil {
		return nil, err
	}
	return credentials, nil
}

func RemoveCredential(ctx *gin.Context, credentialID uuid.UUID, caller uuid.UUID) error {
	if err := VerifyCredentialManageAccessForUser(ctx, credentialID, caller); err != nil {
		return err
	}
	err := repository.RemoveCredential(ctx, credentialID)
	if err != nil {
		return err
	}
	return nil
}
