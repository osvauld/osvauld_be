package controllers

import (
	"errors"
	"net/http"
	"osvauld/customerrors"
	dto "osvauld/dtos"
	"osvauld/infra/logger"
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

		if _, ok := err.(*customerrors.UserDoesNotHaveCredentialAccessError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", err)
			return
		}

		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}
	SendResponse(ctx, http.StatusOK, users, "fetched credential users", nil)
}

func GetCredentialUsersForDataSync(ctx *gin.Context) {

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

	users, err := service.GetCredentialUsersForDataSync(ctx, credentialID, caller)
	if err != nil {

		if _, ok := err.(*customerrors.UserDoesNotHaveCredentialAccessError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", err)
			return
		}

		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}
	SendResponse(ctx, http.StatusOK, users, "fetched credential users", nil)
}

func GetCredentialGroups(ctx *gin.Context) {

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

	groups, err := service.GetCredentialGroups(ctx, credentialID, caller)
	logger.Debugf("Groups: %v", groups)
	if err != nil {

		if _, ok := err.(*customerrors.UserDoesNotHaveCredentialAccessError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", err)
			return
		}

		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}
	SendResponse(ctx, http.StatusOK, groups, "fetched credential groups", nil)
}
