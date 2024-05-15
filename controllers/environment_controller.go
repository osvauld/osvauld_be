package controllers

import (
	"errors"
	"net/http"
	dto "osvauld/dtos"
	"osvauld/service"
	"osvauld/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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

func GetEnvironmentFields(ctx *gin.Context) {
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
	credentials, err := service.GetEnvironmentFields(ctx, environmentID)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}
	SendResponse(ctx, http.StatusOK, credentials, "fetched credentials", nil)
}

func GetEnvironmentByName(ctx *gin.Context) {
	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}
	environmentName := ctx.Param("name")
	environment, err := service.GetEnvironmentByName(ctx, environmentName, caller)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, nil, "", errors.New("failed to fetch environment"))
		return
	}
	SendResponse(ctx, http.StatusOK, environment, "Fetched environment", nil)
}
