package repository

import (
	db "osvauld/db/sqlc"
	"osvauld/infra/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func FetchFieldDetailsByFieldIDForUser(ctx *gin.Context, fieldID uuid.UUID, userID uuid.UUID) (db.FetchFieldNameAndTypeByFieldIDForUserRow, error) {

	fieldDetails, err := database.Store.FetchFieldNameAndTypeByFieldIDForUser(ctx, db.FetchFieldNameAndTypeByFieldIDForUserParams{
		ID:     fieldID,
		UserID: userID,
	})
	if err != nil {
		return fieldDetails, err
	}
	return fieldDetails, nil
}
