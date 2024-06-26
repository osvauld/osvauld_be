package repository

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddEnvironment(ctx *gin.Context, args db.AddEnvironmentParams) (uuid.UUID, error) {
	return database.Store.AddEnvironment(ctx, args)
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
func GetEnvironmentFields(ctx *gin.Context, envID uuid.UUID) ([]db.GetEnvFieldsRow, error) {
	return database.Store.GetEnvFields(ctx, envID)
}

func GetEnvironmentFieldsByName(ctx *gin.Context, name string) ([]db.GetEnvironmentFieldsByNameRow, error) {
	return database.Store.GetEnvironmentFieldsByName(ctx, name)
}

func EditEnvironmentFieldName(ctx *gin.Context, args db.EditEnvironmentFieldNameByIDParams) (string, error) {
	return database.Store.EditEnvironmentFieldNameByID(ctx, args)
}

func IsEnvironmentOwner(ctx *gin.Context, args db.IsEnvironmentOwnerParams) (bool, error) {
	return database.Store.IsEnvironmentOwner(ctx, args)
}

func GetEnvFieldsForCredential(ctx *gin.Context, credentialID uuid.UUID) ([]db.GetEnvFieldsForCredentialRow, error) {
	return database.Store.GetEnvFieldsForCredential(ctx, credentialID)
}

func GetEnvForCredential(ctx *gin.Context, credentialId uuid.UUID) ([]db.GetEnvForCredentialRow, error) {
	return database.Store.GetEnvForCredential(ctx, credentialId)
}
