package controllers

import (
	"errors"
	dto "osvauld/dtos"
	"osvauld/infra/logger"
	service "osvauld/service"
	"osvauld/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateFolder(ctx *gin.Context) {
	var req dto.CreateFolderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Errorf(err.Error())
		SendResponse(ctx, 400, nil, "", err)
	}

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, 401, nil, "", errors.New("unauthorized"))
		return
	}

	folderDetails, err := service.CreateFolder(ctx, req, caller)
	if err != nil {
		logger.Errorf(err.Error())
		SendResponse(ctx, 500, nil, "Failed to create folder", errors.New("failed to create folder"))
		return
	}

	SendResponse(ctx, 200, folderDetails, "", nil)
}

func FetchAccessibleFoldersForUser(ctx *gin.Context) {

	userID, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, 401, nil, "", errors.New("unauthorized"))
		return
	}

	folders, err := service.FetchAccessibleFoldersForUser(ctx, userID)
	if err != nil {
		logger.Errorf(err.Error())
		SendResponse(ctx, 500, nil, "", errors.New("failed to fetch required folders"))
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
	groups, err := service.GetSharedGroupsForFolder(ctx, folderID)
	if err != nil {
		SendResponse(ctx, 500, nil, "", errors.New("failed to fetch required groups"))
		return
	}

	SendResponse(ctx, 200, groups, "Fetched Groups", nil)
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
