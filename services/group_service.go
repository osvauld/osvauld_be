package service

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddGroup(ctx *gin.Context, group dto.CreateGroup, userID uuid.UUID) error {
	err := repository.AddGroup(ctx, group, userID)
	return err
}

func AddMembersToGroup(ctx *gin.Context, payload dto.AddMembers, userID uuid.UUID) error {

	err := repository.AddMembersToGroup(ctx, payload, userID)
	return err
}

func GetUserGroups(ctx *gin.Context, userID uuid.UUID) ([]db.Grouping, error) {
	groups, err := repository.GetUserGroups(ctx, userID)
	return groups, err

}

func GetGroupMembers(ctx *gin.Context, userID uuid.UUID, groupId uuid.UUID) ([]db.GetGroupMembersRow, error) {
	users, err := repository.GetGroupMembers(ctx, groupId)
	return users, err
}
