package repository

import (
	db "osvauld/db/sqlc"
	"osvauld/infra/database"

	"github.com/gin-gonic/gin"
)

func ShareCredentials(ctx *gin.Context, args db.ShareCredentialTransactionParams) error {

	err := database.Store.ShareCredentialsTransaction(ctx, args)
	if err != nil {
		return err
	}

	return nil
}
