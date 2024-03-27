package repository

import (
	"database/sql"
	db "osvauld/db/sqlc"
	"osvauld/infra/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateUser(ctx *gin.Context, args db.CreateUserParams) (uuid.UUID, error) {
	return database.Store.CreateUser(ctx, args)
}

func GetAllUsers(ctx *gin.Context) ([]db.GetAllUsersRow, error) {
	return database.Store.GetAllUsers(ctx)
}

func CreateChallenge(ctx *gin.Context, args db.CreateChallengeParams) (db.SessionTable, error) {
	return database.Store.CreateChallenge(ctx, args)
}

func GetUserByPubKey(ctx *gin.Context, pubKey sql.NullString) (uuid.UUID, error) {
	return database.Store.GetUserByPublicKey(ctx, pubKey)
}

func FetchChallenge(ctx *gin.Context, userId uuid.UUID) (string, error) {
	return database.Store.FetchChallenge(ctx, userId)
}

func UpdateKeys(ctx *gin.Context, args db.UpdateKeysParams) error {
	return database.Store.UpdateKeys(ctx, args)
}

func GetTempPassword(ctx *gin.Context, userName string) (db.GetUserTempPasswordRow, error) {
	return database.Store.GetUserTempPassword(ctx, userName)
}

func UpdateRegistrationChallenge(ctx *gin.Context, args db.UpdateRegistrationChallengeParams) error {
	return database.Store.UpdateRegistrationChallenge(ctx, args)
}

func GetRegistrationChallenge(ctx *gin.Context, userName string) (db.GetRegistrationChallengeRow, error) {
	return database.Store.GetRegistrationChallenge(ctx, userName)
}

func CheckAnyUserExists(ctx *gin.Context) (bool, error) {
	return database.Store.CheckIfUsersExist(ctx)
}

func RemoveUserFromAll(ctx *gin.Context, userID uuid.UUID) error {
	return database.Store.RemoveUserFromAll(ctx, userID)
}
