package repository

import (
	db "osvauld/db/sqlc"
	"osvauld/infra/database"
	"osvauld/infra/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddToAccessList(ctx *gin.Context, credentialID uuid.UUID, accessType string, userID uuid.UUID) error {
	arg := db.AddToAccessListParams{
		CredentialID: credentialID,
		UserID:       userID,
	}
	_, err := database.Store.AddToAccessList(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	return nil
}


func GetUsersByCredential(ctx *gin.Context, credentailID uuid.UUID) ([]db.GetUsersByCredentialRow, error) {
	users, err := database.Store.GetUsersByCredential(ctx, credentailID)
	if err != nil {
		logger.Errorf(err.Error())
		return users, err
	}
	return users, nil
}
