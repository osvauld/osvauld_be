package repository

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func FetchEncryptedFieldsByCredentialIDByAndUserID(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) ([]dto.Field, error) {
	arg := db.FetchEncryptedFieldsByCredentialIDAndUserIDParams{
		CredentialID: credentialID,
		UserID:       userID,
	}
	encryptedData, err := database.Store.FetchEncryptedFieldsByCredentialIDAndUserID(ctx, arg)
	if err != nil {
		return nil, err
	}

	encryptedFields := []dto.Field{}

	for _, encryptedField := range encryptedData {
		encryptedFields = append(
			encryptedFields,
			dto.Field{
				ID:         encryptedField.ID,
				FieldName:  encryptedField.FieldName,
				FieldValue: encryptedField.FieldValue,
			},
		)
	}

	return encryptedFields, err
}
