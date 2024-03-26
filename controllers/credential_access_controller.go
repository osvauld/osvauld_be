package controllers

import (
	"errors"
	"net/http"
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

		if _, ok := err.(*customerrors.UserNotManagerOfCredentialError); ok {
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

		if _, ok := err.(*customerrors.UserNotManagerOfCredentialError); ok {
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

		if _, ok := err.(*customerrors.UserNotManagerOfCredentialError); ok {
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

		if _, ok := err.(*customerrors.UserNotManagerOfFolderError); ok {
			SendResponse(ctx, 401, nil, "", err)
			return
		}

		SendResponse(ctx, 500, nil, "", err)
		return
	}
	SendResponse(ctx, 200, nil, "removed folder access for group", nil)
}

// controller to edit credential access
func EditCredentialAccessForUser(ctx *gin.Context) {
	var req dto.EditCredentialAccessForUser
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

	err = service.EditCredentialAccessForUser(ctx, credentialID, req, caller)
	if err != nil {

		if _, ok := err.(*customerrors.UserNotManagerOfCredentialError); ok {
			SendResponse(ctx, 401, nil, "", err)
			return
		}

		SendResponse(ctx, 500, nil, "", err)
		return
	}
	SendResponse(ctx, 200, nil, "edited credential access", nil)
}

// controller to edit folder access
func EditFolderAccessForUser(ctx *gin.Context) {
	var req dto.EditFolderAccessForUser
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

	err = service.EditFolderAccessForUser(ctx, folderID, req, caller)
	if err != nil {

		if _, ok := err.(*customerrors.UserNotManagerOfFolderError); ok {
			SendResponse(ctx, 401, nil, "", err)
			return
		}

		SendResponse(ctx, 500, nil, "", err)
		return
	}
	SendResponse(ctx, 200, nil, "edited folder access", nil)
}

// controller for edit credential access for user
func EditCredentialAccessForGroup(ctx *gin.Context) {
	var req dto.EditCredentialAccessForGroup
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

	err = service.EditCredentialAccessForGroup(ctx, credentialID, req, caller)
	if err != nil {

		if _, ok := err.(*customerrors.UserNotManagerOfCredentialError); ok {
			SendResponse(ctx, 401, nil, "", err)
			return
		}

		SendResponse(ctx, 500, nil, "", err)
		return
	}
	SendResponse(ctx, 200, nil, "edited credential access for group", nil)
}

// controller for edit folder access for group
func EditFolderAccessForGroup(ctx *gin.Context) {
	var req dto.EditFolderAccessForGroup
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

	err = service.EditFolderAccessForGroup(ctx, folderID, req, caller)
	if err != nil {

		if _, ok := err.(*customerrors.UserNotManagerOfFolderError); ok {
			SendResponse(ctx, 401, nil, "", err)
			return
		}

		SendResponse(ctx, 500, nil, "", err)
		return
	}
	SendResponse(ctx, 200, nil, "edited folder access for group", nil)
}

func GetCredentialUsersWithDirectAccess(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusUnauthorized, nil, "Unauthorized", err)
		return
	}

	credentialIDStr := ctx.Param("id")
	credentialID, err := uuid.Parse(credentialIDStr)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid credential id"))
		return
	}

	users, err := service.GetCredentialUsersWithDirectAccess(ctx, credentialID, caller)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}
	SendResponse(ctx, http.StatusOK, users, "fetched credential users", nil)
}

func GetCredentialUsersWithAllAccessSource(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusUnauthorized, nil, "Unauthorized", err)
		return
	}

	credentialIDStr := ctx.Param("id")
	credentialID, err := uuid.Parse(credentialIDStr)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid credential id"))
		return
	}

	users, err := service.GetCredentialUsersWithAllAccessSource(ctx, credentialID, caller)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}
	SendResponse(ctx, http.StatusOK, users, "fetched credential users", nil)
}
