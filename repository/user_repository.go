package repository

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/database"
	"osvauld/infra/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateUser(ctx *gin.Context, user dto.CreateUser) (uuid.UUID, error) {
	arg := db.CreateUserParams{
		Username:  user.UserName,
		Name:      user.Name,
		PublicKey: user.PublicKey,
	}
	id, err := database.Q.CreateUser(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return uuid.Nil, err
	}
	return id, err
}

func GetUser(ctx *gin.Context, userLogin dto.Login) (db.User, error) {
	user, err := database.Q.GetUserByUsername(ctx, userLogin.UserName)
	if err != nil {
		logger.Errorf(err.Error())
		return user, err
	}
	return user, nil
}
