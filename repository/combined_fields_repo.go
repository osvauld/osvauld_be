package repository

import (
	"osvauld/infra/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GenerateCombinedFieldEntry(ctx *gin.Context, userID uuid.UUID) error {

	return database.Store.CreateCombinedField(ctx, userID)
}
