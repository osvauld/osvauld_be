package service

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetEnvironmentFields(ctx *gin.Context, envID uuid.UUID) ([]db.GetEnvFieldsRow, error) {
	envFields, err := repository.GetEnvironmentFields(ctx, envID)
	if err != nil {
		return []db.GetEnvFieldsRow{}, err
	}
	return envFields, nil
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
