package service

import (
	db "osvauld/db/sqlc"
	"osvauld/infra/database"
	"osvauld/infra/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAccessForCredential(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) (bool, error) {
	arg := db.HasUserAccessParams{
		CredentialID: credentialID,
		UserID:       userID,
	}
	hasAccess, err := database.Store.HasUserAccess(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return hasAccess, err
	}
	return hasAccess, nil
}
