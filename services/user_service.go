package service

import (
	"fmt"
	"osvauld/auth"
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateUser(ctx *gin.Context, user dto.CreateUser) (uuid.UUID, error) {
	id, err := repository.CreateUser(ctx, user)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func GetAllUsers(ctx *gin.Context) ([]db.GetAllUsersRow, error) {
	users, err := repository.GetAllUsers(ctx)
	if err != nil {
		return users, err
	}
	return users, nil
}
func Login(ctx *gin.Context, userData dto.Login) dto.LoginReturn {
	user, _ := repository.GetUser(ctx, userData)
	token, _ := auth.GenerateToken(user.Username, user.ID)
	fmt.Println(token)
	loginReturn := dto.LoginReturn{
		User:  user,
		Token: token,
	}

	return loginReturn
}