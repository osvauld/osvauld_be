package controllers

import (
	"errors"
	"net/http"
	dto "osvauld/dtos"
	service "osvauld/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddGroup(ctx *gin.Context) {
	var req dto.CreateGroup
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIdInterface, _ := ctx.Get("userId")
	userID, _ := userIdInterface.(uuid.UUID)
	err := service.AddGroup(ctx, req, userID)
	if err != nil {
		SendResponse(ctx, 500, nil, "failed to create group", errors.New("failed to create user"))
		return
	}
	SendResponse(ctx, 201, nil, "created group", nil)

}

func AppendMembersToGroup(ctx *gin.Context) {
	var req dto.AddMembers

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIdInterface, _ := ctx.Get("userId")
	userID, _ := userIdInterface.(uuid.UUID)

	// Check user is autorized to append users to group
	isMember, err := service.CheckUserMemberOfGroup(ctx, userID, req.GroupID)
	if !isMember {
		SendResponse(ctx, 401, nil, "", errors.New("unauthorized for group append"))
		return
	}
	if err != nil {
		SendResponse(ctx, 500, nil, "Failed to append user", errors.New("failed to append user"))
		return
	}

	err = service.AddMembersToGroup(ctx, req, userID)
	if err != nil {
		SendResponse(ctx, 500, nil, "Failed to append user", errors.New("failed to append user"))
		return
	}
	SendResponse(ctx, 201, nil, "Added users to group", nil)

}

func GetUserGroups(ctx *gin.Context) {
	userIdInterface, _ := ctx.Get("userId")
	userID, _ := userIdInterface.(uuid.UUID)
	groups, err := service.GetUserGroups(ctx, userID)
	if err != nil {
		SendResponse(ctx, 500, nil, "Failed to fetch Groups", errors.New("failed to fetch users"))
		return
	}
	SendResponse(ctx, 200, groups, "Fetched user groups", nil)

}

func GetGroupMembers(ctx *gin.Context) {

	userIdInterface, _ := ctx.Get("userId")
	userID, _ := userIdInterface.(uuid.UUID)

	groupIDStr := ctx.Param("groupId")
	groupID, _ := uuid.Parse(groupIDStr)

	// Check user is autorized to see members of the group
	isMember, err := service.CheckUserMemberOfGroup(ctx, userID, groupID)
	if !isMember {
		SendResponse(ctx, 401, nil, "", errors.New("unauthorized for group view"))
		return
	}
	if err != nil {
		SendResponse(ctx, 500, nil, "Failed to fetch Group members", errors.New("failed to fetch group members"))
		return
	}

	groups, err := service.GetGroupMembers(ctx, userID, groupID)
	if err != nil {
		SendResponse(ctx, 500, nil, "Failed to fetch Group members", errors.New("failed to fetch group members"))
		return
	}
	SendResponse(ctx, 200, groups, "Fetched group memebers", nil)
}
