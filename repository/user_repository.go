package repository

import (
	dto "osvauld/dtos"
	"osvauld/infra/database"
	"osvauld/infra/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateUser(ctx *gin.Context, user dto.CreateUser) (uuid.UUID, error) {
	id, err := database.Q.CreateUser(ctx, user.UserName)
	if err != nil {
		logger.Errorf(err.Error())
		return uuid.Nil, err
	}
	return id, err
}
