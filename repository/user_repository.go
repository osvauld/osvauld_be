package repository

import (
	"database/sql"
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/database"
	"osvauld/infra/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateUser(ctx *gin.Context, user dto.CreateUser) (uuid.UUID, error) {
	arg := db.CreateUserParams{
		Username:     user.UserName,
		Name:         user.Name,
		TempPassword: user.TempPassword,
	}
	id, err := database.Store.CreateUser(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return uuid.Nil, err
	}
	return id, err
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
	nullPubKey := sql.NullString{String: pubKey, Valid: true}
	user, err := database.Store.GetUserByPublicKey(ctx, nullPubKey)
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

func FetchChallenge(ctx *gin.Context, userId uuid.UUID) (string, error) {
	challenge, err := database.Store.FetchChallenge(ctx, userId)
	if err != nil {
		logger.Errorf(err.Error())
		return challenge, err
	}
	return challenge, nil
}

func CheckTempPassword(ctx *gin.Context, password string, username string) (bool, error) {
	arg := db.CheckTempPasswordParams{
		Username:     username,
		TempPassword: password,
	}
	count, err := database.Store.CheckTempPassword(ctx, arg)
	if err != nil || count == 0 {
		logger.Errorf(err.Error())
		return false, err
	}
	return true, nil
}

func UpdateKeys(ctx *gin.Context, username string, encryptionKey string, deviceKey string) error {
	encryptionKeyNull := sql.NullString{String: encryptionKey, Valid: true}
	deviceKeyNull := sql.NullString{String: deviceKey, Valid: true}
	arg := db.UpdateKeysParams{
		Username:      username,
		EncryptionKey: encryptionKeyNull,
		DeviceKey:     deviceKeyNull,
	}
	err := database.Store.UpdateKeys(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	return nil
}
