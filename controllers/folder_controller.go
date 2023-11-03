package controllers

import (
	"net/http"
	dto "osvauld/dtos"
	"osvauld/infra/logger"
	service "osvauld/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateFolder(ctx *gin.Context) {
	var req dto.CreateFolder
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Errorf(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse request"})
	}

	userIdString := ctx.GetHeader("userId")
	if userIdString == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id header is required"})
		return
	}
	userID, _ := uuid.Parse(userIdString)
	// TODO add validation
	service.CreateFolder(ctx, req, userID)
}

func GetAccessibleFolders(ctx *gin.Context) {
	// 1. Fetch User ID from Header
	userIDString := ctx.GetHeader("userId")
	if userIDString == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID not provided"})
		return
	}
	userID, _ := uuid.Parse(userIDString)
	folders, _ := service.GetAccessibleFolders(ctx, userID)
	// 5. Return the Folders to the Client
	ctx.JSON(http.StatusOK, gin.H{"folders": folders})
}
