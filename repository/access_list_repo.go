package repository

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/database"
	"osvauld/infra/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddToAccessList(ctx *gin.Context, credentialID uuid.UUID, accessType string, userID uuid.UUID) error {
	arg := db.AddToAccessListParams{
		CredentialID: credentialID,
		UserID:       userID,
	}
	_, err := database.Store.AddToAccessList(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	return nil
}

func GetUsersByCredential(ctx *gin.Context, credentailID uuid.UUID) ([]db.GetUsersByCredentialRow, error) {
	users, err := database.Store.GetUsersByCredential(ctx, credentailID)
	if err != nil {
		logger.Errorf(err.Error())
		return users, err
	}
	return users, nil
}

func GetCredentialAccessForUser(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) ([]dto.AccessListResult, error) {
	arg := db.GetCredentialAccessForUserParams{
		CredentialID: credentialID,
		UserID:       userID,
	}

	access, err := database.Store.GetCredentialAccessForUser(ctx, arg)
	if err != nil {
		return []dto.AccessListResult{}, err
	}
	accessListResults := []dto.AccessListResult{}
	for _, access := range access {
		accessListResults = append(accessListResults, dto.AccessListResult{
			ID:         access.ID,
			UserID:     access.UserID,
			AccessType: access.AccessType,
			GroupID:    access.GroupID,
		})
	}
	return accessListResults, nil
}
