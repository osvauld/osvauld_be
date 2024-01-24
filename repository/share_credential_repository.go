package repository

import (
	dto "osvauld/dtos"
	"osvauld/infra/database"

	"github.com/gin-gonic/gin"
)

func ShareCredentials(ctx *gin.Context, args dto.ShareCredentialTransactionParams) error {

	err := database.Store.ShareCredentialsTransaction(ctx, args)
	if err != nil {
		return err
	}

	return nil
}
