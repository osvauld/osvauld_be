package repository

import (
	db "osvauld/db/sqlc"
	"osvauld/infra/database"

	"github.com/gin-gonic/gin"
)

func GetAllFieldsForCredentialIDs(ctx *gin.Context, args db.GetAllFieldsForCredentialIDsParams) ([]db.GetAllFieldsForCredentialIDsRow, error) {
	return database.Store.GetAllFieldsForCredentialIDs(ctx, args)
}

func GetNonSensitiveFieldsForCredentialIDs(ctx *gin.Context, args db.GetNonSensitiveFieldsForCredentialIDsParams) ([]db.GetNonSensitiveFieldsForCredentialIDsRow, error) {
	return database.Store.GetNonSensitiveFieldsForCredentialIDs(ctx, args)
}

func CheckFieldEntryExists(ctx *gin.Context, args db.CheckFieldEntryExistsParams) (bool, error) {
	return database.Store.CheckFieldEntryExists(ctx, args)
}

func GetSensitiveFields(ctx *gin.Context, args db.GetSensitiveFieldsParams) ([]db.GetSensitiveFieldsRow, error) {
	return database.Store.GetSensitiveFields(ctx, args)
}

func DeleteAccessRemovedFields(ctx *gin.Context) error {
	return database.Store.DeleteAccessRemovedFields(ctx)
}
