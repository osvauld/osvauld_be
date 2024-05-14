package controllers

import (
	"errors"
	"net/http"
	"osvauld/customerrors"
	dto "osvauld/dtos"
	"osvauld/infra/logger"
	"osvauld/service"
	"osvauld/utils"

	"github.com/gin-gonic/gin"
)

func ShareCredentialsWithUsers(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusUnauthorized, nil, "", errors.New("unauthorized"))
		return
	}

	var req dto.ShareCredentialsWithUsersRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", err)
		return
	}

	response, err := service.ShareCredentialsWithUsers(ctx, req.UserData, caller)
	if err != nil {
		if _, ok := err.(*customerrors.UserNotManagerOfCredentialError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", err)
			return
		}
		logger.Errorf(err.Error())
		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}
	SendResponse(ctx, http.StatusOK, response, "Success", nil)
}

func ShareCredentialsWithGroups(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusUnauthorized, nil, "Unauthorized", errors.New("unauthorized"))
		return
	}

	var req dto.ShareCredentialsWithGroupsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", err)
		return
	}

	response, err := service.ShareCredentialsWithGroups(ctx, req.GroupData, caller)
	if err != nil {
		if _, ok := err.(*customerrors.UserNotManagerOfCredentialError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", err)
			return
		}
		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}
	SendResponse(ctx, http.StatusOK, response, "Success", nil)
}

func ShareFolderWithUsers(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusUnauthorized, nil, "", errors.New("unauthorized"))
		return
	}

	var req dto.ShareFolderWithUsersRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", err)
		return
	}

	responses, err := service.ShareFolderWithUsers(ctx, req, caller)
	if err != nil {
		if _, ok := err.(*customerrors.UserNotManagerOfCredentialError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", err)
			return
		}
		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}
	SendResponse(ctx, http.StatusOK, responses, "Success", nil)
}

func ShareFolderWithGroups(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusUnauthorized, nil, "", errors.New("unauthorized"))
		return
	}

	var req dto.ShareFolderWithGroupsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", err)
		return
	}

	response, err := service.ShareFolderWithGroups(ctx, req, caller)
	if err != nil {
		if _, ok := err.(*customerrors.UserNotManagerOfCredentialError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", err)
			return
		}
		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}
	SendResponse(ctx, http.StatusOK, response, "Success", nil)
}

func ShareCredentialsWithEnvironment(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusUnauthorized, nil, "", errors.New("unauthorized"))
		return
	}

	var req dto.ShareCredentialsWithEnvironmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", err)
		return
	}

	err = service.ShareCredentialsWithEnvironment(ctx, req, caller)
	if err != nil {

		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}
	SendResponse(ctx, http.StatusOK, nil, "Success", nil)

}
