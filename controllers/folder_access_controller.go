package controllers

import (
	"errors"
	"net/http"
	"osvauld/customerrors"
	dto "osvauld/dtos"
	"osvauld/service"
	"osvauld/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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

func GetFolderUsersForDataSync(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}

	folderIDStr := ctx.Param("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid folder id"))
		return
	}

	users, err := service.GetFolderUsersForDataSync(ctx, folderID, caller)
	if err != nil {

		if _, ok := err.(*customerrors.UserDoesNotHaveFolderAccessError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", err)
			return
		}

		SendResponse(ctx, http.StatusInternalServerError, nil, "", errors.New("failed to fetch required users"))
		return
	}

	SendResponse(ctx, http.StatusOK, users, "Fetched users", nil)
}

func GetFolderUsersWithDirectAccess(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}

	folderIDStr := ctx.Param("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid folder id"))
		return
	}

	users, err := service.GetFolderUsersWithDirectAccess(ctx, folderID, caller)
	if err != nil {

		if _, ok := err.(*customerrors.UserDoesNotHaveFolderAccessError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", err)
			return
		}

		SendResponse(ctx, http.StatusInternalServerError, nil, "", errors.New("failed to fetch required users"))
		return
	}

	SendResponse(ctx, http.StatusOK, users, "Fetched users", nil)
}

func GetFolderGroups(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}

	folderIDStr := ctx.Param("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid folder id"))
		return
	}

	groups, err := service.GetFolderGroups(ctx, folderID, caller)
	if err != nil {

		if _, ok := err.(*customerrors.UserDoesNotHaveFolderAccessError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", err)
			return
		}

		SendResponse(ctx, http.StatusInternalServerError, nil, "", errors.New("failed to fetch required groups"))
		return
	}

	SendResponse(ctx, http.StatusOK, groups, "Fetched groups", nil)
}
