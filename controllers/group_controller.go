package controllers

import (
	"fmt"
	"net/http"
	"osvauld/infra/logger"
	"osvauld/models"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateGroupRequest struct {
	Name string `json:"name" binding:"required"`
}

func AddGroup(ctx *gin.Context) {
	var req CreateGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extracting the user from the header
	userIDHeader := ctx.GetHeader("user_id")
	if userIDHeader == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "UserID header missing"})
		return
	}

	userID, err := uuid.Parse(userIDHeader)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UserID format"})
		return
	}

	group := models.Group{
		Name:      req.Name,
		Members:   []uuid.UUID{userID},
		CreatedBy: userID,
	}
	err = repository.AddGroup(group)
	if err != nil {
		logger.Errorf(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to create group"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Group created successfully", "group": group})
}

type UpdateGroupMembersRequest struct {
	GroupID string   `json:"group_id"`
	UserIDs []string `json:"user_ids"`
}

func AppendMembersToGroup(ctx *gin.Context) {
	var req UpdateGroupMembersRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse UUIDs
	var userUUIDs []uuid.UUID
	for _, idStr := range req.UserIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid UUID format for UserID: %s", idStr)})
			return
		}
		userUUIDs = append(userUUIDs, id)
	}

	groupID, err := uuid.Parse(req.GroupID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format for GroupID"})
		return
	}

	err = repository.AddMembersToGroup(groupID, userUUIDs)
	if err != nil {
		logger.Errorf(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to add members to group"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Members added successfully"})
}
