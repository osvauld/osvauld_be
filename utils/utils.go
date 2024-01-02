package utils

import (
	"errors"

	"crypto/rand"

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

func CreateRandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic(err) // Handle the error properly in production code
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes)
}
