package controllers

import (
	"errors"
	"net/http"
	"osvauld/customerrors"
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
		SendResponse(ctx, http.StatusBadRequest, nil, "", err)
	}

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}

	folderDetails, err := service.CreateFolder(ctx, req, caller)
	if err != nil {
		logger.Errorf(err.Error())
		SendResponse(ctx, http.StatusInternalServerError, nil, "", errors.New("failed to create folder"))
		return
	}

	SendResponse(ctx, http.StatusOK, folderDetails, "", nil)
}

func FetchAccessibleFoldersForUser(ctx *gin.Context) {

	userID, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}

	folders, err := service.FetchAccessibleFoldersForUser(ctx, userID)
	if err != nil {
		logger.Errorf(err.Error())
		SendResponse(ctx, http.StatusInternalServerError, nil, "", errors.New("failed to fetch required folders"))
		return
	}

	SendResponse(ctx, http.StatusOK, folders, "Fetched folders", nil)

}

func GetSharedUsersForFolder(ctx *gin.Context) {

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

	users, err := service.GetSharedUsersForFolder(ctx, folderID, caller)
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

func GetSharedGroupsForFolder(ctx *gin.Context) {

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

	groups, err := service.GetSharedGroupsForFolder(ctx, folderID, caller)
	if err != nil {

		if _, ok := err.(*customerrors.UserDoesNotHaveFolderAccessError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", err)
			return
		}

		SendResponse(ctx, http.StatusInternalServerError, nil, "", errors.New("failed to fetch required groups"))
		return
	}

	SendResponse(ctx, http.StatusOK, groups, "Fetched Groups", nil)
}

func GetGroupsWithoutAccess(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}

	folderIDStr := ctx.Param("folderId")
	folderID, err := uuid.Parse(folderIDStr)

	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid folder id"))
		return
	}

	groups, err := service.GetGroupsWithoutAccess(ctx, folderID, caller)
	if err != nil {
		SendResponse(ctx, 500, nil, "Failed to get groups", errors.New("failed to fetch required groups"))
		return
	}

	SendResponse(ctx, 200, groups, "Fetched groups", nil)
}
