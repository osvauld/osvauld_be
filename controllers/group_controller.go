package controllers

import (
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

	userIdString := ctx.GetHeader("userId")
	if userIdString == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id header is required"})
		return
	}
	userID, _ := uuid.Parse(userIdString)
	service.AddGroup(ctx, req, userID)
}

func AppendMembersToGroup(ctx *gin.Context) {
	var req dto.AddMembers

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIdString := ctx.GetHeader("userId")
	if userIdString == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id header is required"})
		return
	}
	userID, _ := uuid.Parse(userIdString)
	service.AddMembersToGroup(ctx, req, userID)
}
