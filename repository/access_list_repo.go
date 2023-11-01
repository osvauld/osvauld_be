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
		AccessType:   db.AccessType(accessType),
		UserID:       uuid.NullUUID{UUID: userID, Valid: true},
	}
	q := db.New(database.DB)
	_, err := q.AddToAccessList(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	return nil
}
