package controllers

import (
	"errors"
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

	userIdInterface, _ := ctx.Get("userId")
	userID, _ := userIdInterface.(uuid.UUID)
	err := service.CreateFolder(ctx, req, userID)
	if err != nil {
		SendResponse(ctx, 500, nil, "Failed to create folder", errors.New("failed to create folder"))
		return
	}
	SendResponse(ctx, 201, nil, "Created folder", nil)
}

func GetAccessibleFolders(ctx *gin.Context) {
	// 1. Fetch User ID from Header
	userIdInterface, _ := ctx.Get("userId")
	userID, _ := userIdInterface.(uuid.UUID)
	folders, err := service.GetAccessibleFolders(ctx, userID)
	if err != nil {
		SendResponse(ctx, 500, nil, "Failed to get folders", errors.New("failed to fetch required folders"))
		return
	}

	SendResponse(ctx, 200, folders, "Fetched folders", nil)

}
