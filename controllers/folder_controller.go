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

func RemoveFolder(ctx *gin.Context) {
	folderIDStr := ctx.Param("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid folder id"))
		return
	}
	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}
	err = service.RemoveFolder(ctx, folderID, caller)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, nil, "", errors.New("failed to remove folder"))
		return
	}

	SendResponse(ctx, http.StatusOK, nil, "Folder removed successfully", nil)
}

func EditFolder(ctx *gin.Context) {
	var payload dto.EditFolder
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", err)
	}

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}

	folderID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid folder id"))
		return
	}

	err = service.EditFolder(ctx, folderID, payload, caller)
	if err != nil {

		if _, ok := err.(*customerrors.UserNotManagerOfFolderError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", err)
			return
		}

		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}

	SendResponse(ctx, http.StatusOK, nil, "Folder edited successfully", nil)
}

func AddEnvironment(ctx *gin.Context) {
	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}
	var req dto.AddEnvironment
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", err)
		return
	}
	_, err = service.AddEnvironment(ctx, req, caller)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}
	SendResponse(ctx, http.StatusOK, nil, "added environment", nil)
}

func GetEnvironments(ctx *gin.Context) {
	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}
	environments, err := service.GetEnvironments(ctx, caller)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}
	SendResponse(ctx, http.StatusOK, environments, "fetched environments", nil)
}

func GetEnvironmentCredentials(ctx *gin.Context) {
	// caller, err := utils.FetchUserIDFromCtx(ctx)
	// if err != nil {
	// 	SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
	// 	return
	// }
	environmentIDStr := ctx.Param("id")
	environmentID, err := uuid.Parse(environmentIDStr)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid environment id"))
		return
	}
	// TODO: Add check for user access to environment
	credentials, err := service.GetEnvironmentCredentials(ctx, environmentID)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}
	SendResponse(ctx, http.StatusOK, credentials, "fetched credentials", nil)
}
