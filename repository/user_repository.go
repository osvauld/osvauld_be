package repository

import (
	"database/sql"
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateUser(ctx *gin.Context, args db.CreateUserParams) (uuid.UUID, error) {
	return database.Store.CreateUser(ctx, args)
}

func GetAllSignedUpUsers(ctx *gin.Context) ([]db.GetAllSignedUpUsersRow, error) {
	return database.Store.GetAllSignedUpUsers(ctx)
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
	return database.Store.RemoveUserFromOrg(ctx, userID)
}

func CheckNameExists(ctx *gin.Context, name string) (bool, error) {
	return database.Store.CheckNameExist(ctx, name)
}

func CheckUsernameExists(ctx *gin.Context, username string) (bool, error) {
	return database.Store.CheckUsernameExist(ctx, username)

}

func CheckUserType(ctx *gin.Context, userID uuid.UUID) (string, error) {
	return database.Store.GetUserType(ctx, userID)
}

func GetUserByID(ctx *gin.Context, userID uuid.UUID) (db.GetUserByIDRow, error) {
	return database.Store.GetUserByID(ctx, userID)
}

func GetAllUsers(ctx *gin.Context) ([]db.GetAllUsersRow, error) {
	return database.Store.GetAllUsers(ctx)
}

func GetUserDeviceKey(ctx *gin.Context, userID uuid.UUID) (string, error) {
	return database.Store.GetUserDeviceKey(ctx, userID)
}

func CreateCLIUser(ctx *gin.Context, userDetails dto.CreateCLIUser, caller uuid.UUID) (uuid.UUID, error) {
	return database.Store.CreateCliUser(ctx, db.CreateCliUserParams{
		Name:          userDetails.Name,
		DeviceKey:     sql.NullString{String: userDetails.DeviceKey, Valid: true},
		EncryptionKey: sql.NullString{String: userDetails.EncryptionKey, Valid: true},
		Username:      userDetails.Name,
		CreatedBy:     uuid.NullUUID{UUID: caller, Valid: true},
		TempPassword:  userDetails.Name,
		Type:          "cli",
		Status:        "active",
		SignedUp:      true,
	})
}

func GetEnvironments(ctx *gin.Context, userID uuid.UUID) ([]db.GetEnvironmentsForUserRow, error) {
	return database.Store.GetEnvironmentsForUser(ctx, uuid.NullUUID{UUID: userID, Valid: true})
}

func GetCliUsers(ctx *gin.Context, userID uuid.UUID) ([]db.GetCliUsersRow, error) {
	return database.Store.GetCliUsers(ctx, uuid.NullUUID{UUID: userID, Valid: true})
}

func GetSuperUser(ctx *gin.Context) (db.User, error) {
	return database.Store.GetSuperUser(ctx)
}
