package service

import (
	"errors"
	"osvauld/auth"
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/logger"
	"osvauld/repository"
	"osvauld/utils"

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

// func Login(ctx *gin.Context, userData dto.Login) (dto.LoginReturn, error) {
// 	user, _ := repository.GetUser(ctx, userData)
// 	// if err != nil {
// 	// 	return dto.LoginReturn{}, err
// 	// }
// 	token, _ := auth.GenerateToken(user.Username, user.ID)
// 	fmt.Println(token)
// 	loginReturn := dto.LoginReturn{
// 		User:  user.ID.String(),
// 		Token: token,
// 	}

// 	return loginReturn, nil
// }

func CreateChallenge(ctx *gin.Context, user dto.CreateChallenge) (string, error) {
	challengeStr := utils.CreateRandomString(12)
	logger.Debugf("challenge string: %s", challengeStr)
	userId, err := repository.GetUserByPubKey(ctx, user.PublicKey)
	if err != nil || userId == uuid.Nil {
		return "", err
	}
	challenge, err := repository.CreateChallenge(ctx, user.PublicKey, challengeStr, userId)
	if err != nil {
		return challenge.Challenge, err
	}
	return challenge.Challenge, nil
}

func VerifyChallenge(ctx *gin.Context, challenge dto.VerifyChallenge) (string, error) {
	userId, err := repository.GetUserByPubKey(ctx, challenge.PublicKey)
	if err != nil || userId == uuid.Nil {
		logger.Errorf(err.Error())
		return "", err
	}
	challengeStr, err := repository.FetchChallenge(ctx, userId)
	if err != nil {
		logger.Errorf(err.Error())
		return "", err
	}
	resp, err := auth.VerifySignature(challenge.Signature, challenge.PublicKey, challengeStr, userId)
	if err != nil || userId == uuid.Nil {
		return "", err
	}
	return resp, nil

}

func Register(ctx *gin.Context, registerData dto.Register) (bool, error) {
	pass, err := repository.CheckTempPassword(ctx, registerData.Password, registerData.UserName)
	if err != nil {
		logger.Errorf(err.Error())
		return false, err
	}
	if !pass {
		logger.Errorf("password not matched")
		return false, errors.New("password not matched")
	}

	err = repository.UpdateKeys(ctx, registerData.UserName, registerData.RsaKey, registerData.EccKey)
	if err != nil {
		logger.Errorf(err.Error())
		return false, err
	}
	return true, nil

}
