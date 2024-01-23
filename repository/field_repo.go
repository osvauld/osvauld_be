package repository

import (
	db "osvauld/db/sqlc"
	"osvauld/infra/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func FetchFieldDetailsByFieldID(ctx *gin.Context, fieldID uuid.UUID) (db.FetchFieldNameAndTypeByFieldIDRow, error) {
	fieldDetails, err := database.Store.FetchFieldNameAndTypeByFieldID(ctx, fieldID)
	if err != nil {
		return fieldDetails, err
	}
	return fieldDetails, nil
}
