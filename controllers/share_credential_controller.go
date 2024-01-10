package controllers

import (
	"errors"
	"net/http"
	dto "osvauld/dtos"
	"osvauld/service"
	"osvauld/utils"

	"github.com/gin-gonic/gin"
)

func ShareMultipleCredentialsWithMulitpleUsers(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, 401, nil, "Unauthorized", errors.New("unauthorized"))
		return
	}

	var req dto.ShareMultipleCredentialsWithMultipleUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := service.ShareMultipleCredentialsWithMultipleUsers(ctx, req.UserData, caller)
	if err != nil {
		SendResponse(ctx, 500, nil, "Failed to share credential", errors.New("failed to share credential"))
		return
	}
	SendResponse(ctx, 200, response, "Success", nil)
}

func ShareMultipleCredentialsWithMulitpleGroups(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, 401, nil, "Unauthorized", errors.New("unauthorized"))
		return
	}

	var req dto.ShareMultipleCredentialsWithMultipleGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := service.ShareMultipleCredentialsWithMulitpleGroups(ctx, req.GroupData, caller)
	if err != nil {
		SendResponse(ctx, 500, nil, "Failed to share credential", errors.New("failed to share credential"))
		return
	}
	SendResponse(ctx, 200, response, "Success", nil)
}
