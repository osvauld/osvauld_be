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

// controller to remove folder access for user
func RemoveFolderAccessForUsers(ctx *gin.Context) {
	var req dto.RemoveFolderAccessForUsers
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendResponse(ctx, 400, nil, "", err)
		return
	}

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, 401, nil, "Unauthorized", err)
		return
	}

	folderIDStr := ctx.Param("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		SendResponse(ctx, 400, nil, "", err)
		return
	}

	err = service.RemoveFolderAccessForUsers(ctx, folderID, req, caller)
	if err != nil {

		if _, ok := err.(*customerrors.UserNotAnOwnerOfFolderError); ok {
			SendResponse(ctx, 401, nil, "", err)
			return
		}

		SendResponse(ctx, 500, nil, "", err)
		return
	}
	SendResponse(ctx, 200, nil, "removed folder access for user", nil)
}

// controller to remove credential access for group
func RemoveCredentialAccessForGroups(ctx *gin.Context) {
	var req dto.RemoveCredentialAccessForGroups
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

	err = service.RemoveCredentialAccessForGroups(ctx, credentialID, req, caller)
	if err != nil {

		if _, ok := err.(*customerrors.UserNotAnOwnerOfCredentialError); ok {
			SendResponse(ctx, 401, nil, "", err)
			return
		}

		SendResponse(ctx, 500, nil, "", err)
		return
	}
	SendResponse(ctx, 200, nil, "removed credential access for group", nil)
}

// controller to remove folder access for group
func RemoveFolderAccessForGroups(ctx *gin.Context) {
	var req dto.RemoveFolderAccessForGroups
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendResponse(ctx, 400, nil, "", err)
		return
	}

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, 401, nil, "Unauthorized", err)
		return
	}

	folderIDStr := ctx.Param("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		SendResponse(ctx, 400, nil, "", err)
		return
	}

	err = service.RemoveFolderAccessForGroups(ctx, folderID, req, caller)
	if err != nil {

		if _, ok := err.(*customerrors.UserNotAnOwnerOfFolderError); ok {
			SendResponse(ctx, 401, nil, "", err)
			return
		}

		SendResponse(ctx, 500, nil, "", err)
		return
	}
	SendResponse(ctx, 200, nil, "removed folder access for group", nil)
}
