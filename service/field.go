package service

import (
	dto "osvauld/dtos"
	"osvauld/infra/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func FetchFieldNameAndTypeByFieldID(ctx *gin.Context, fieldID uuid.UUID) (dto.Field, error) {
	fieldDetails, err := database.Store.FetchFieldNameAndTypeByFieldID(ctx, fieldID)
	if err != nil {
		return dto.Field{}, err
	}

	fieldDto := dto.Field{
		ID:        fieldID,
		FieldName: fieldDetails.FieldName,
		FieldType: fieldDetails.FieldType,
	}

	return fieldDto, nil
}
