package controllers

import (
	"osvauld/customerrors"
	dto "osvauld/dtos"
	"osvauld/utils"

	"osvauld/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// controller for remove credential access for user
func RemoveCredentialAccessForUsers(ctx *gin.Context) {
	var req dto.RemoveCredentialAccessForUsers
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendResponse(ctx, 400, nil, "", err)
		return
	}

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, 401, nil, "Unauthorized", err)
		return
	}

	credentialIDStr := ctx.Param("id")
	credentialID, err := uuid.Parse(credentialIDStr)
	if err != nil {
		SendResponse(ctx, 400, nil, "", err)
		return
	}

	err = service.RemoveCredentialAccessForUsers(ctx, credentialID, req, caller)
	if err != nil {

		if _, ok := err.(*customerrors.UserNotAnOwnerOfCredentialError); ok {
			SendResponse(ctx, 401, nil, "", err)
			return
		}

		SendResponse(ctx, 500, nil, "", err)
		return
	}
	SendResponse(ctx, 200, nil, "removed credential access for user", nil)
}
