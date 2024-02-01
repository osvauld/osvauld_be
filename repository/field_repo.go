package repository

import (
	db "osvauld/db/sqlc"
	"osvauld/infra/database"

	"github.com/gin-gonic/gin"
)

func GetFieldDataForCredentialIDsForUser(ctx *gin.Context, args db.GetFieldDataByCredentialIDsForUserParams) ([]db.GetFieldDataByCredentialIDsForUserRow, error) {

	return database.Store.GetFieldDataByCredentialIDsForUser(ctx, args)
}

func CheckFieldEntryExists(ctx *gin.Context, args db.CheckFieldEntryExistsParams) (bool, error) {

	return database.Store.CheckFieldEntryExists(ctx, args)
}
