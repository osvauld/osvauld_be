package repository

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/database"
	"osvauld/infra/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SaveEncryptedData(ctx *gin.Context, encrypedData dto.FieldRequest, credentialID uuid.UUID, userID uuid.UUID) (uuid.UUID, error) {
	arg := db.CreateEncryptedDataParams{
		FieldName:    encrypedData.FieldName,
		FieldValue:   encrypedData.FieldValue,
		UserID:       uuid.NullUUID{UUID: userID, Valid: true},
		CredentialID: uuid.NullUUID{UUID: credentialID, Valid: true},
	}
	q := db.New(database.DB)
	id, err := q.CreateEncryptedData(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return uuid.Nil, err
	}
	return id, nil
}
