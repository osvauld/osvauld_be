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

	hashedPassword, err := auth.HashPassword(user.TempPassword)
	if err != nil {
		return uuid.Nil, err
	}

	return repository.CreateUser(ctx, db.CreateUserParams{
		Username:     user.UserName,
		Name:         user.Name,
		TempPassword: hashedPassword,
	})

}

func TempLogin(ctx *gin.Context, req dto.TempLogin) (string, error) {
	user, err := repository.GetTempPassword(ctx, req.UserName)
	if err != nil {
		logger.Errorf(err.Error())
		return "", errors.New("user not found")
	}

	if user.Status != "created" {
		return "", errors.New("temp login not allowed")
	}

	tempPasswordHash := user.TempPassword
	passwordMatched := auth.CheckPasswordHash(req.TempPassword, tempPasswordHash)
	if !passwordMatched {
		return "", errors.New("incorrect password")
	}

	challengeStr := utils.CreateRandomString(12)

	err = repository.UpdateRegistrationChallenge(ctx, db.UpdateRegistrationChallengeParams{
		Username:              req.UserName,
		RegistrationChallenge: sql.NullString{String: challengeStr, Valid: true},
	})
	if err != nil {
		return "", err
	}

	return challengeStr, nil
}

func Register(ctx *gin.Context, registerData dto.Register) (bool, error) {

	registrationChallenge, err := repository.GetRegistrationChallenge(ctx, registerData.UserName)
	if err != nil {
		return false, err
	}

	if registrationChallenge.Status != "temp_login" {
		return false, errors.New("registration not allowed")
	}

	valid, err := auth.VerifySignature(registerData.Signature, registerData.DeviceKey, registrationChallenge.RegistrationChallenge.String)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, errors.New("invalid signature")
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
	if err != nil {
		logger.Errorf(err.Error())
		return "", err
	}
	challengeStr, err := repository.FetchChallenge(ctx, userId)
	if err != nil {
		logger.Errorf(err.Error())
		return "", err
	}
	valid, err := auth.VerifySignature(challenge.Signature, challenge.PublicKey, challengeStr)
	if err != nil {
		return "", err
	}

	if !valid {
		return "", errors.New("invalid signature")
	}

	token, err := auth.GenerateToken("test", userId)
	if err != nil {
		return "", err
	}

	return token, nil

}

func CheckUserExists(ctx *gin.Context) (bool, error) {
	check, err := repository.CheckAnyUserExists(ctx)
	if err != nil {
		return false, err
	}
	return check, nil
}

func RemoveUserFromAll(ctx *gin.Context, userID uuid.UUID) error {
	return repository.RemoveUserFromAll(ctx, userID)
}
