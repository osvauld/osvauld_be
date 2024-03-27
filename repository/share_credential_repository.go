package repository

import (
	db "osvauld/db/sqlc"
	"osvauld/infra/database"

	"github.com/gin-gonic/gin"
)

func ShareCredentials(ctx *gin.Context, args db.ShareCredentialTransactionParams) error {
	return database.Store.ShareCredentialsTransaction(ctx, args)
}
