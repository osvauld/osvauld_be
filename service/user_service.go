package service

import (
	"database/sql"
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
	return repository.CreateUser(ctx, db.CreateUserParams{
		Username:     user.UserName,
		Name:         user.Name,
		TempPassword: user.TempPassword,
	})

}

func GetAllUsers(ctx *gin.Context) ([]db.GetAllUsersRow, error) {

	return repository.GetAllUsers(ctx)
}

func CreateChallenge(ctx *gin.Context, user dto.CreateChallenge) (string, error) {
	challengeStr := utils.CreateRandomString(12)

	userId, err := repository.GetUserByPubKey(ctx, sql.NullString{String: user.PublicKey, Valid: true})
	if err != nil || userId == uuid.Nil {
		return "", err
	}

	challenge, err := repository.CreateChallenge(ctx, db.CreateChallengeParams{
		UserID:    userId,
		PublicKey: user.PublicKey,
		Challenge: challengeStr,
	})
	if err != nil {
		return challenge.Challenge, err
	}
	return challenge.Challenge, nil
}

func VerifyChallenge(ctx *gin.Context, challenge dto.VerifyChallenge) (string, error) {
	userId, err := repository.GetUserByPubKey(ctx, sql.NullString{String: challenge.PublicKey, Valid: true})
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

	pass, err := repository.CheckTempPassword(ctx, db.CheckTempPasswordParams{
		Username:     registerData.UserName,
		TempPassword: registerData.Password,
	})
	if err != nil {
		logger.Errorf(err.Error())
		return false, err
	}
	if !pass {
		logger.Errorf("password not matched")
		return false, errors.New("password not matched")
	}

	err = repository.UpdateKeys(ctx, db.UpdateKeysParams{
		Username:      registerData.UserName,
		EncryptionKey: sql.NullString{String: registerData.EncryptionKey, Valid: true},
		DeviceKey:     sql.NullString{String: registerData.DeviceKey, Valid: true},
	})

	if err != nil {
		logger.Errorf(err.Error())
		return false, err
	}
	return true, nil

}

func GetCredentialUsers(ctx *gin.Context, credentialID uuid.UUID) ([]db.GetAccessTypeAndUsersByCredentialIdRow, error) {

	return repository.GetCredentialUsers(ctx, credentialID)
}
