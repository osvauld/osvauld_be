package repository

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddEnvironment(ctx *gin.Context, args dto.AddEnvironment, caller uuid.UUID) (uuid.UUID, error) {
	return database.Store.AddEnvironment(ctx, db.AddEnvironmentParams{
		Name:      args.Name,
		CliUser:   args.CliUser,
		CreatedBy: caller,
	})
}

func CheckCredentialExistsInEnvironment(ctx *gin.Context, credentialID uuid.UUID, environmentID uuid.UUID) (bool, error) {
	return database.Store.CheckCredentialExistsForEnv(ctx, db.CheckCredentialExistsForEnvParams{
		CredentialID: credentialID,
		EnvID:        environmentID,
	})
}

func AddCredentialFieldsToEnvironment(ctx *gin.Context, args []dto.CredentialEnvData) error {
	return database.Store.AddCredentialFieldToEnvTxn(ctx, args)
}

func GetEnvironmentByID(ctx *gin.Context, environmentID uuid.UUID, caller uuid.UUID) (db.Environment, error) {
	return database.Store.GetEnvironmentByID(ctx, db.GetEnvironmentByIDParams{
		ID:        environmentID,
		CreatedBy: caller,
	})
}
func GetEnvironmentCredentials(ctx *gin.Context, envID uuid.UUID) ([]db.GetEnvFieldsRow, error) {
	return database.Store.GetEnvFields(ctx, envID)
}
