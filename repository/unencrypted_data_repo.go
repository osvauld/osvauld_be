package repository

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/database"
	"osvauld/infra/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SaveUnEncryptedData(ctx *gin.Context, encrypedData dto.FieldRequest, credentiaId uuid.UUID) (uuid.UUID, error) {
	arg := db.CreateUnencryptedDataParams{
		FieldName:    encrypedData.FieldName,
		FieldValue:   encrypedData.FieldValue,
		CredentialID: uuid.NullUUID{UUID: credentiaId, Valid: true},
	}
	q := db.New(database.DB)
	id, err := q.CreateUnencryptedData(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return uuid.Nil, err
	}
	return id, nil
}
