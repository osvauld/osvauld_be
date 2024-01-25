package controllers

import (
	"errors"
	"net/http"
	dto "osvauld/dtos"
	"osvauld/infra/logger"
	"osvauld/service"
	"osvauld/utils"

	"github.com/gin-gonic/gin"
)

func ShareCredentialsWithUsers(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, 401, nil, "Unauthorized", errors.New("unauthorized"))
		return
	}

	var req dto.ShareCredentialsWithUsersRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := service.ShareCredentialsWithUsers(ctx, req.UserData, caller)
	if err != nil {
		logger.Errorf(err.Error())
		SendResponse(ctx, 500, nil, "Failed to share credential", errors.New("failed to share credential"))
		return
	}
	SendResponse(ctx, 200, response, "Success", nil)
}

func ShareCredentialsWithGroups(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, 401, nil, "Unauthorized", errors.New("unauthorized"))
		return
	}

	var req dto.ShareCredentialsWithGroupsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := service.ShareCredentialsWithGroups(ctx, req.GroupData, caller)
	if err != nil {
		SendResponse(ctx, 500, nil, "Failed to share credential", errors.New("failed to share credential"))
		return
	}
	SendResponse(ctx, 200, response, "Success", nil)
}

func ShareFolderWithUsers(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, 401, nil, "Unauthorized", errors.New("unauthorized"))
		return
	}

	var req dto.ShareFolderWithUsersRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	responses, err := service.ShareFolderWithUsers(ctx, req, caller)
	if err != nil {
		SendResponse(ctx, 400, nil, "", errors.New("failed to share"))
		return
	}
	SendResponse(ctx, 200, responses, "Success", nil)
}

func ShareFolderWithGroups(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, 401, nil, "Unauthorized", errors.New("unauthorized"))
		return
	}

	var req dto.ShareFolderWithGroupsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := service.ShareFolderWithGroups(ctx, req, caller)
	if err != nil {
		SendResponse(ctx, 400, nil, "Failed to share folder with groups", nil)
		return
	}
	SendResponse(ctx, 200, response, "Success", nil)
}
