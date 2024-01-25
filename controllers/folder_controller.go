package controllers

import (
	"errors"
	"net/http"
	dto "osvauld/dtos"
	"osvauld/infra/logger"
	service "osvauld/service"

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

	userIdInterface, _ := ctx.Get("userId")
	userID, _ := userIdInterface.(uuid.UUID)
	folders, err := service.GetAccessibleFolders(ctx, userID)
	if err != nil {
		SendResponse(ctx, 500, nil, "Failed to get folders", errors.New("failed to fetch required folders"))
		return
	}

	SendResponse(ctx, 200, folders, "Fetched folders", nil)

}

func GetUsersByFolder(ctx *gin.Context) {

	userIdInterface, _ := ctx.Get("userId")
	userID, _ := userIdInterface.(uuid.UUID)

	folderIDStr := ctx.Param("id")
	folderID, _ := uuid.Parse(folderIDStr)
	users, err := service.GetUsersByFolder(ctx, folderID, userID)

	if err != nil {
		SendResponse(ctx, 500, nil, "Failed to get users", errors.New("failed to fetch required users"))
		return
	}

	SendResponse(ctx, 200, users, "Fetched users", nil)
}

func GetSharedUsersForFolder(ctx *gin.Context) {
	folderIDStr := ctx.Param("id")
	folderID, _ := uuid.Parse(folderIDStr)
	users, err := service.GetSharedUsersForFolder(ctx, folderID)
	if err != nil {
		SendResponse(ctx, 500, nil, "Failed to get users", errors.New("failed to fetch required users"))
		return
	}

	SendResponse(ctx, 200, users, "Fetched users", nil)
}

func GetSharedGroupsForFolder(ctx *gin.Context) {
	folderIDStr := ctx.Param("id")
	folderID, _ := uuid.Parse(folderIDStr)
	groups, err := service.GetSharedUsersForFolder(ctx, folderID)
	if err != nil {
		SendResponse(ctx, 500, nil, "", errors.New("failed to fetch required groups"))
		return
	}

	SendResponse(ctx, 200, groups, "Fetched users", nil)
}

func GetGroupsWithoutAccess(ctx *gin.Context) {
	folderIDStr := ctx.Param("folderId")
	folderID, _ := uuid.Parse(folderIDStr)
	groups, err := service.GetGroupsWithoutAccess(ctx, folderID)
	if err != nil {
		SendResponse(ctx, 500, nil, "Failed to get groups", errors.New("failed to fetch required groups"))
		return
	}

	SendResponse(ctx, 200, groups, "Fetched groups", nil)
}
