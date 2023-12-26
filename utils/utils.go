package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func FetchUserIDFromCtx(ctx *gin.Context) (uuid.UUID, error) {
	userIdInterface, _ := ctx.Get("userId")
	userID, ok := userIdInterface.(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("failed to fetch user id from context")
	}
	return userID, nil
}
