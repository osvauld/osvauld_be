package service

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetSensitiveFields(ctx *gin.Context, credentialID uuid.UUID, caller uuid.UUID) ([]dto.Field, error) {

	if err := VerifyCredentialReadAccessForUser(ctx, credentialID, caller); err != nil {
		return nil, err
	}

	sensitiveFields, err := repository.GetSensitiveFields(ctx, db.GetSensitiveFieldsParams{
		CredentialID: credentialID,
		UserID:       caller,
	})

	fieldObjs := []dto.Field{}

	for _, field := range sensitiveFields {
		fieldObjs = append(fieldObjs, dto.Field{
			ID:         field.ID,
			FieldName:  field.FieldName,
			FieldValue: field.FieldValue,
			FieldType:  "sensitive",
		})
	}

	return fieldObjs, err
}

func GetCredentialsFieldsByFolderID(ctx *gin.Context, folderID uuid.UUID, caller uuid.UUID) ([]dto.CredentialFields, error) {

	if err := VerifyFolderReadAccessForUser(ctx, folderID, caller); err != nil {
		return nil, err
	}

	// we could merge the below two queries into one, but it is two queries not a big deal
	credentialsIds, err := repository.GetCredentialIdsByFolderAndUserId(ctx, db.GetCredentialIdsByFolderParams{
		FolderID: folderID,
		UserID:   caller,
	})
	if err != nil {
		return nil, err
	}

	credentials, err := GetFieldsByCredentialIDs(ctx, credentialsIds, caller)
	if err != nil {
		return nil, err
	}

	return credentials, nil
}

func GetFieldsByCredentialIDs(ctx *gin.Context, credentialIDs []uuid.UUID, caller uuid.UUID) ([]dto.CredentialFields, error) {

	if err := VerifyReadAccessForCredentials(ctx, credentialIDs, caller); err != nil {
		return nil, err
	}

	fields, err := repository.GetAllFieldsForCredentialIDs(ctx, db.GetAllFieldsForCredentialIDsParams{
		Credentials: credentialIDs,
		UserID:      caller,
	})

	if err != nil {
		return nil, err
	}

	credentialMap := make(map[uuid.UUID][]dto.Field)

	for _, field := range fields {

		fieldObj := dto.Field{
			ID:         field.ID,
			FieldName:  field.FieldName,
			FieldValue: field.FieldValue,
			FieldType:  field.FieldType,
		}

		credentialMap[field.CredentialID] = append(credentialMap[field.CredentialID], fieldObj)

	}

	credentialFieldDtos := []dto.CredentialFields{}

	for _, credentialID := range credentialIDs {
		credentialFieldDtos = append(credentialFieldDtos, dto.CredentialFields{
			CredentialID: credentialID,
			Fields:       credentialMap[credentialID],
		})
	}

	return credentialFieldDtos, nil
}

func DeleteAccessRemovedFields(ctx *gin.Context) error {
	return repository.DeleteAccessRemovedFields(ctx)
}
