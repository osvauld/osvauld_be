package service

import (
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
