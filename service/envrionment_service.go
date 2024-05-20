package service

import (
	"osvauld/customerrors"
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func VerifyEnvironmentAccessForUser(ctx *gin.Context, environmentID uuid.UUID, userID uuid.UUID) error {

	hasAccess, err := repository.IsEnvironmentOwner(ctx, db.IsEnvironmentOwnerParams{
		ID:        environmentID,
		CreatedBy: userID,
	})
	if err != nil {
		return err
	}
	if !hasAccess {
		return &customerrors.UserDoesNotHaveEnvironmentAccess{UserID: userID, EnvironmentID: environmentID}
	}

	return nil
}

func GetEnvironmentFields(ctx *gin.Context, envID uuid.UUID) ([]dto.CredentialEnvFields, error) {
	envFields, err := repository.GetEnvironmentFields(ctx, envID)
	credentialFieldMap := make(map[uuid.UUID][]dto.EnvFieldData)
	credentialIDNameMap := make(map[uuid.UUID]string)
	for _, field := range envFields {
		credentialFieldMap[field.CredentialID] = append(credentialFieldMap[field.CredentialID], dto.EnvFieldData{
			FieldID:    field.ID,
			FieldName:  field.FieldName,
			FieldValue: field.FieldValue,
		})
		credentialIDNameMap[field.CredentialID] = field.CredentialName
	}

	if err != nil {
		return []dto.CredentialEnvFields{}, err
	}

	var credentialEnvData = []dto.CredentialEnvFields{}

	for credentialID, fields := range credentialFieldMap {
		credentialEnvData = append(credentialEnvData, dto.CredentialEnvFields{
			CredentialID: credentialID,
			CredentialName: credentialIDNameMap[credentialID],
			Fields:       fields,
		})
	}
	return credentialEnvData, nil
}

func GetEnvironmentByName(ctx *gin.Context, environmentName string, userID uuid.UUID) ([]db.GetEnvironmentFieldsByNameRow, error) {
	// TODO: validation
	environment, err := repository.GetEnvironmentFieldsByName(ctx, environmentName)
	if err != nil {
		return []db.GetEnvironmentFieldsByNameRow{}, err
	}
	return environment, nil
}

func AddEnvironment(ctx *gin.Context, environment dto.AddEnvironment, caller uuid.UUID) (uuid.UUID, error) {
	// TODO: verify no duplicate name for user
	return repository.AddEnvironment(ctx, db.AddEnvironmentParams{
		Name:      environment.Name,
		CliUser:   environment.CliUser,
		CreatedBy: caller,
	})
}

func GetEnvironments(ctx *gin.Context, userID uuid.UUID) ([]db.GetEnvironmentsForUserRow, error) {
	return repository.GetEnvironments(ctx, userID)
}

func EditEnvFieldName(ctx *gin.Context, payload dto.EditEnvFieldName, caller uuid.UUID) (map[string]string, error) {

	if err := VerifyEnvironmentAccessForUser(ctx, payload.EnvironmentID, caller); err != nil {
		return nil, err
	}

	fieldName, err := repository.EditEnvironmentFieldName(ctx, db.EditEnvironmentFieldNameByIDParams{
		ID:        payload.FieldID,
		FieldName: payload.FieldName,
		EnvID:     payload.EnvironmentID,
	})
	if err != nil {
		return nil, err
	}

	return map[string]string{"fieldName": fieldName}, nil
}
