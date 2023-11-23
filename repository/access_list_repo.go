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
		CredentialID: uuid.NullUUID{UUID: credentialID, Valid: true},
		AccessType:   accessType,
		UserID:       uuid.NullUUID{UUID: userID, Valid: true},
	}
	_, err := database.Q.AddToAccessList(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	return nil
}

func CheckAccessForCredential(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) (bool, error) {
	arg := db.HasUserAccessParams{
		CredentialID: uuid.NullUUID{UUID: credentialID, Valid: true},
		UserID:       uuid.NullUUID{UUID: userID, Valid: true},
	}
	hasAccess, err := database.Q.HasUserAccess(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return hasAccess, err
	}
	return hasAccess, nil
}

func GetUsersByCredential(ctx *gin.Context, credentailID uuid.UUID) ([]db.GetUsersByCredentialRow, error) {
	users, err := database.Q.GetUsersByCredential(ctx, uuid.NullUUID{UUID: credentailID, Valid: true})
	if err != nil {
		logger.Errorf(err.Error())
		return users, err
	}
	return users, nil
}
