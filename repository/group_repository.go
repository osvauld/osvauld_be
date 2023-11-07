package repository

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/database"
	"osvauld/infra/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func AddGroup(ctx *gin.Context, group dto.CreateGroup, userID uuid.UUID) error {
	arg := db.CreateGroupParams{
		Name:      group.Name,
		CreatedBy: uuid.NullUUID{UUID: userID, Valid: true},
		Members:   []uuid.UUID{userID},
	}
	_, err := database.Q.CreateGroup(ctx, arg)
	return err
}

func AddMembersToGroup(ctx *gin.Context, payload dto.AddMembers, userID uuid.UUID) error {
	uuidArray := pq.Array(payload.Members)
	arg := db.AddMemberToGroupParams{
		CreatedBy: uuid.NullUUID{UUID: userID, Valid: true},
		ID:        payload.GroupID,
		ArrayCat:  uuidArray,
	}
	err := database.Q.AddMemberToGroup(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	return nil

}

func GetUserGroups(ctx *gin.Context, userID uuid.UUID) ([]db.Group, error) {

	groups, err := database.Q.GetUserGroups(ctx, []uuid.UUID{userID})
	if err != nil {
		logger.Errorf(err.Error())
		return groups, err
	}
	return groups, nil
}

func GetGroupMembers(ctx *gin.Context, groupID uuid.UUID) ([]db.GetGroupMembersRow, error) {
	users, err := database.Q.GetGroupMembers(ctx, []uuid.UUID{groupID})
	if err != nil {
		logger.Errorf(err.Error())
		return users, err
	}
	return users, nil
}
