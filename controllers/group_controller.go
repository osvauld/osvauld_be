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

func CreateGroup(ctx *gin.Context) {
	var req dto.CreateGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", err)
		return
	}

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}

	groupDetails, err := service.AddGroup(ctx, req.Name, caller)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, nil, "", errors.New("failed to create group"))
		return
	}

	SendResponse(ctx, http.StatusCreated, groupDetails, "created group", nil)

}

func GetUserGroups(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}

	groups, err := service.GetUserGroups(ctx, caller)
	if err != nil {

		SendResponse(ctx, http.StatusInternalServerError, nil, "", errors.New("failed to fetch users"))
		return
	}
	SendResponse(ctx, http.StatusOK, groups, "Fetched user groups", nil)

}

func GetGroupMembers(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("unauthorized"))
		return
	}

	groupIDStr := ctx.Param("groupId")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid group id"))
		return
	}

	groups, err := service.GetGroupMembers(ctx, groupID, caller)
	if err != nil {
		if _, ok := err.(*customerrors.UserNotMemberOfGroupError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", err)
			return
		}

		SendResponse(ctx, http.StatusInternalServerError, nil, "", errors.New("failed to fetch group members"))
		return
	}
	SendResponse(ctx, http.StatusOK, groups, "Fetched group memebers", nil)
}

func GetAllCredentialFieldsByGroupID(ctx *gin.Context) {

	userID, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, 401, nil, "", errors.New("unauthorized"))
		return
	}

	groupIDStr := ctx.Param("groupId")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		SendResponse(ctx, 400, nil, "", errors.New("invalid group id"))
		return
	}

	credentialFields, err := service.GetCredentialFieldsByGroupID(ctx, userID, groupID)
	if err != nil {
		if _, ok := err.(*customerrors.UserNotMemberOfGroupError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", err)
			return
		}
		SendResponse(ctx, http.StatusInternalServerError, nil, "", errors.New("failed to fetch credential"))
		return
	}

	SendResponse(ctx, http.StatusOK, credentialFields, "Fetched credentials", nil)
}

func AddMemberToGroup(ctx *gin.Context) {

	var req dto.AddMemberToGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", err)
		return
	}

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}

	err = service.AddMemberToGroup(ctx, req, caller)
	if err != nil {
		logger.Errorf(err.Error())
		if _, ok := err.(*customerrors.UserAlreadyMemberOfGroupError); ok {
			SendResponse(ctx, http.StatusConflict, nil, "", err)
			return
		} else if _, ok := err.(*customerrors.UserNotAdminOfGroupError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", err)
			return
		}
		SendResponse(ctx, http.StatusInternalServerError, nil, "failed to add members to group", nil)
		return
	}

	SendResponse(ctx, http.StatusOK, nil, "Added members to group", nil)
}

func GetUsersOfGroups(ctx *gin.Context) {

	var req dto.GetUsersOfGroupsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", err)
		return
	}

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}

	groupUsers, err := service.GetUsersOfGroups(ctx, req.GroupIDs, caller)
	if err != nil {

		if _, ok := err.(*customerrors.UserNotMemberOfGroupError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", err)
			return
		}

		SendResponse(ctx, http.StatusInternalServerError, nil, "", errors.New("failed to fetch group users"))
		return
	}
	SendResponse(ctx, 200, groupUsers, "Fetched group users", nil)
}

func GetCredentialGroups(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}

	credentialIDStr := ctx.Param("id")
	credentialID, err := uuid.Parse(credentialIDStr)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid credential id"))
		return
	}

	users, err := service.GetCredentialGroups(ctx, credentialID, caller)
	if err != nil {
		if _, ok := err.(*customerrors.UserDoesNotHaveCredentialAccessError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", err)
			return
		}
		SendResponse(ctx, http.StatusInternalServerError, nil, "failed to get credential users", err)
		return
	}
	SendResponse(ctx, http.StatusOK, users, "fetched credential users", nil)
}

func GetUsersWithoutGroupAccess(ctx *gin.Context) {
	userID, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}

	groupIDStr := ctx.Param("groupId")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid group id"))
		return
	}

	users, err := service.GetUsersWithoutGroupAccess(ctx, groupID, userID)
	if err != nil {
		if _, ok := err.(*customerrors.UserNotMemberOfGroupError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", err)
			return
		}
		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}
	SendResponse(ctx, http.StatusOK, users, "Fetched users not in group", nil)
}
