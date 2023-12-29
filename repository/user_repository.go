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
	id, err := database.Store.CreateUser(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return uuid.Nil, err
	}
	return id, err
}

func GetUser(ctx *gin.Context, userLogin dto.Login) (db.GetUserByUsernameRow, error) {
	user, err := database.Store.GetUserByUsername(ctx, userLogin.UserName)
	if err != nil {
		logger.Errorf(err.Error())
		return user, err
	}
	return user, nil
}

func GetAllUsers(ctx *gin.Context) ([]db.GetAllUsersRow, error) {
	user, err := database.Store.GetAllUsers(ctx)
	if err != nil {
		logger.Errorf(err.Error())
		return user, err
	}
	return user, nil
}

func GetUserByPubKey(ctx *gin.Context, pubKey string) (uuid.UUID, error) {
	user, err := database.Store.GetUserByPublicKey(ctx, pubKey)
	if err != nil {
		logger.Errorf(err.Error())
		return user, err
	}
	return user, nil
}

func CreateChallenge(ctx *gin.Context, pubKey string, challenge string, userId uuid.UUID) (db.SessionTable, error) {
	arg := db.CreateChallengeParams{
		PublicKey: pubKey,
		Challenge: challenge,
		UserID:    userId,
	}
	challengeRow, err := database.Store.CreateChallenge(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return challengeRow, err
	}
	return challengeRow, nil
}
