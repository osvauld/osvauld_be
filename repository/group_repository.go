package repository

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/database"
	"osvauld/infra/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddGroup(ctx *gin.Context, group dto.CreateGroup, userID uuid.UUID) error {
	arg := db.CreateGroupParams{
		Name:   group.Name,
		UserID: userID,
	}
	err := database.Q.CreateGroup(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	return nil
}

func AddMembersToGroup(ctx *gin.Context, payload dto.AddMembers, userID uuid.UUID) error {
	// uuidArray := pq.Array(payload.Members)
	arg := db.AddMemberToGroupParams{
		GroupingID: payload.GroupID,
		Column2:    payload.Members,
	}
	err := database.Q.AddMemberToGroup(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	return nil

}

func GetUserGroups(ctx *gin.Context, userID uuid.UUID) ([]db.Grouping, error) {

	groups, err := database.Q.GetUserGroups(ctx, userID)
	if err != nil {
		logger.Errorf(err.Error())
		return groups, err
	}
	return groups, nil
}

func GetGroupMembers(ctx *gin.Context, groupID uuid.UUID) ([]db.GetGroupMembersRow, error) {
	users, err := database.Q.GetGroupMembers(ctx, groupID)
	if err != nil {
		logger.Errorf(err.Error())
		return users, err
	}
	return users, nil
}

func CheckUserMemberOfGroup(ctx *gin.Context, userID uuid.UUID, groupID uuid.UUID) (bool, error) {
	args := db.CheckUserMemberOfGroupParams{
		UserID:     userID,
		GroupingID: groupID,
	}
	isMember, err := database.Q.CheckUserMemberOfGroup(ctx, args)
	if err != nil {
		logger.Errorf(err.Error())
		return false, err
	}
	return isMember, nil
}