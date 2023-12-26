package repository

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/database"
	"osvauld/infra/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SaveUnEncryptedData(ctx *gin.Context, encrypedData dto.Field, credentiaId uuid.UUID) (uuid.UUID, error) {
	arg := db.CreateUnencryptedDataParams{
		FieldName:    encrypedData.FieldName,
		FieldValue:   encrypedData.FieldValue,
		CredentialID: credentiaId,
	}
	id, err := database.Store.CreateUnencryptedData(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return uuid.Nil, err
	}
	return id, nil
}

func FetchUnencryptedFieldsByCredentialID(ctx *gin.Context, credentialID uuid.UUID) ([]dto.Field, error) {

	unEncryptedData, err := database.Store.FetchUnencryptedFieldsByCredentialID(ctx, credentialID)
	if err != nil {
		return nil, err
	}

	unEncryptedFields := []dto.Field{}

	for _, unEncryptedField := range unEncryptedData {
		unEncryptedFields = append(
			unEncryptedFields,
			dto.Field{
				ID:         unEncryptedField.ID,
				FieldName:  unEncryptedField.FieldName,
				FieldValue: unEncryptedField.FieldValue,
			},
		)
	}

	return unEncryptedFields, err
}
