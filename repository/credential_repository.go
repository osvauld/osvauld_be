package repository

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/database"
	"osvauld/infra/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddCredential(ctx *gin.Context, args dto.AddCredentialDto, caller uuid.UUID) (uuid.UUID, error) {

	credentialID, err := database.Store.AddCredentialTransaction(ctx, args, caller)
	if err != nil {
		return uuid.UUID{}, err
	}
	return credentialID, nil
}

func FetchCredentialByID(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) (dto.CredentialDetails, error) {

	credentialDetails := dto.CredentialDetails{}
	credentialDetails.UserID = userID

	credential, err := database.Store.FetchCredentialDataByID(ctx, credentialID)
	if err != nil {
		return dto.CredentialDetails{}, err
	}

	credentialDetails.CredentialID = credential.ID
	credentialDetails.Name = credential.Name
	credentialDetails.FolderID = credential.FolderID
	credentialDetails.CreatedAt = credential.CreatedAt
	credentialDetails.UpdatedAt = credential.UpdatedAt
	credentialDetails.CreatedBy = credential.CreatedBy
	if credential.Description.Valid {
		credentialDetails.Description = credential.Description.String
	}

	return credentialDetails, nil
}

func GetCredentialsByFolder(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) ([]dto.CredentialForUser, error) {

	args := db.FetchCredentialIdsForUserByFolderIdParams{
		FolderID: folderID,
		UserID:   userID,
	}

	// Users can have access to only some of the credentials in a folder.
	// So check the access_list table to see which credentials the user has access to
	credentialDetails, err := database.Store.FetchCredentialIdsForUserByFolderId(ctx, args)
	if err != nil {
		return []dto.CredentialForUser{}, err
	}

	credentialIDs := []uuid.UUID{}
	for _, credential := range credentialDetails {
		credentialIDs = append(credentialIDs, credential.CredentialID)
	}

	arg := db.FetchCredentialFieldsForUserByCredentialIdsParams{
		Column1: credentialIDs,
		UserID:  userID,
	}
	FieldsData, err := database.Store.FetchCredentialFieldsForUserByCredentialIds(ctx, arg)
	if err != nil {
		return []dto.CredentialForUser{}, err
	}

	credentialFieldGroups := map[uuid.UUID][]dto.Field{}

	for _, credential := range FieldsData {
		// if credential.CredentialID does not exist add it to the map and add the field to the array
		if _, ok := credentialFieldGroups[credential.CredentialID]; ok {
			credentialFieldGroups[credential.CredentialID] = append(credentialFieldGroups[credential.CredentialID], dto.Field{
				ID:         credential.FieldID,
				FieldName:  credential.FieldName,
				FieldValue: credential.FieldValue,
				FieldType:  credential.FieldType,
			})
		} else {
			credentialFieldGroups[credential.CredentialID] = []dto.Field{
				{
					ID:         credential.FieldID,
					FieldName:  credential.FieldName,
					FieldValue: credential.FieldValue,
					FieldType:  credential.FieldType,
				},
			}
		}
	}

	credentials := []dto.CredentialForUser{}
	for _, credential := range credentialDetails {
		credentialForUser := dto.CredentialForUser{}
		credentialForUser.CredentialID = credential.CredentialID
		credentialForUser.Name = credential.Name
		credentialForUser.Description = credential.Description
		credentialForUser.CredentialType = credential.CredentialType
		credentialForUser.FolderID = folderID
		credentialForUser.CreatedAt = credential.CreatedAt
		credentialForUser.UpdatedAt = credential.UpdatedAt
		credentialForUser.CreatedBy = credential.CreatedBy
		credentialForUser.Fields = credentialFieldGroups[credential.CredentialID]
		credentials = append(credentials, credentialForUser)
	}

	return credentials, nil
}

func FetchUnEncryptedData(ctx *gin.Context, credentialID uuid.UUID) ([]db.GetCredentialUnencryptedDataRow, error) {

	encryptedData, err := database.Store.GetCredentialUnencryptedData(ctx, credentialID)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	return encryptedData, err
}

func GetCredentialsFieldsByIds(ctx *gin.Context, credentialIds []uuid.UUID, userID uuid.UUID) ([]db.GetCredentialsFieldsByIdsRow, error) {
	arg := db.GetCredentialsFieldsByIdsParams{
		Column1: credentialIds,
		UserID:  userID,
	}
	encryptedData, err := database.Store.GetCredentialsFieldsByIds(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	return encryptedData, err
}

func GetCredentialsByIDs(ctx *gin.Context, credentialIds []uuid.UUID, userID uuid.UUID) ([]db.GetCredentialDetailsByIdsRow, error) {
	credentials, err := database.Store.GetCredentialDetailsByIds(ctx, db.GetCredentialDetailsByIdsParams{
		UserID:  userID,
		Column1: credentialIds,
	})
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	return credentials, err
}

func GetAllUrlsForUser(ctx *gin.Context, userID uuid.UUID) ([]db.GetAllUrlsForUserRow, error) {
	urls, err := database.Store.GetAllUrlsForUser(ctx, userID)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	return urls, err
}

func GetSensitiveFieldsById(ctx *gin.Context, credentialID uuid.UUID, caller uuid.UUID) ([]db.GetSensitiveFieldsRow, error) {
	// Check if caller has access
	sensitiveFields, err := database.Store.GetSensitiveFields(ctx, db.GetSensitiveFieldsParams{
		CredentialID: credentialID,
		UserID:       caller,
	})

	return sensitiveFields, err
}

func GetCredentialIdsByFolderAndUserId(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) ([]uuid.UUID, error) {
	credentialIds, err := database.Store.GetCredentialIdsByFolder(ctx, db.GetCredentialIdsByFolderParams{
		FolderID: folderID,
		UserID:   userID,
	})
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	return credentialIds, err
}

func GetCredentialUsers(ctx *gin.Context, credentialID uuid.UUID) ([]db.GetAccessTypeAndUsersByCredentialIdRow, error) {
	return database.Store.GetAccessTypeAndUsersByCredentialId(ctx, credentialID)

}

func GetCredentialGroups(ctx *gin.Context, credentialID uuid.UUID) ([]db.GetAccessTypeAndGroupsByCredentialIdRow, error) {
	groups, err := database.Store.GetAccessTypeAndGroupsByCredentialId(ctx, credentialID)
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	return groups, err
}
