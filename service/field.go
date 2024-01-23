package service

import (
	dto "osvauld/dtos"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func FetchFieldNameAndTypeByFieldIDForUser(ctx *gin.Context, fieldID uuid.UUID, userID uuid.UUID) (dto.Field, error) {
	fieldDetails, err := repository.FetchFieldDetailsByFieldIDForUser(ctx, fieldID, userID)
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
