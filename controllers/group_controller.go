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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, 401, nil, "", errors.New("unauthorized"))
		return
	}

	groupDetails, err := service.AddGroup(ctx, req.Name, caller)
	if err != nil {
		SendResponse(ctx, 500, nil, "", errors.New("failed to create user"))
		return
	}

	SendResponse(ctx, 201, groupDetails, "created group", nil)

}

func GetUserGroups(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, 401, nil, "", errors.New("unauthorized"))
		return
	}

	groups, err := service.GetUserGroups(ctx, caller)
	if err != nil {
		SendResponse(ctx, 500, nil, "", errors.New("failed to fetch users"))
		return
	}
	SendResponse(ctx, 200, groups, "Fetched user groups", nil)

}

func GetGroupMembers(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
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

	groups, err := service.GetGroupMembers(ctx, groupID, caller)
	if err != nil {

		if _, ok := err.(*customerrors.UserNotAuthenticatedError); ok {
			SendResponse(ctx, 401, nil, err.Error(), nil)
			return
		}

		SendResponse(ctx, 500, nil, "", errors.New("failed to fetch group members"))
		return
	}
	SendResponse(ctx, 200, groups, "Fetched group memebers", nil)
}

func GetAllCredentialsByGroupID(ctx *gin.Context) {

	userID, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, 401, nil, "", errors.New("unauthorized"))
		return
	}

	groupIDStr := ctx.Param("groupId")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		SendResponse(ctx, 400, nil, "Invalid group id", nil)
		return
	}

	credentialFields, err := service.GetCredentialFieldsByGroupID(ctx, userID, groupID)
	if err != nil {
		if _, ok := err.(*customerrors.UserNotAuthenticatedError); ok {
			SendResponse(ctx, 401, nil, "Unauthorized", errors.New("unauthorized"))
			return
		}
		SendResponse(ctx, 200, nil, "Failed to fetch credential", errors.New("failed to fetch credential"))
		return
	}

	if len(credentialFields) == 0 {
		SendResponse(ctx, 204, nil, "No credentials found", nil)
		return
	}

	SendResponse(ctx, 200, credentialFields, "Fetched credentials", nil)
}

func AddMemberToGroup(ctx *gin.Context) {

	userID, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, 401, nil, "Unauthorized", errors.New("unauthorized"))
		return
	}

	var req dto.AddMemberToGroupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = service.AddMemberToGroup(ctx, req, userID)
	if err != nil {
		logger.Errorf(err.Error())
		if _, ok := err.(*customerrors.UserAlreadyMemberOfGroupError); ok {
			SendResponse(ctx, 409, nil, "", err)
			return
		} else if _, ok := err.(*customerrors.UserNotAuthenticatedError); ok {
			SendResponse(ctx, 401, nil, "", err)
			return
		}
		SendResponse(ctx, 500, nil, "failed to add members to group", nil)
		return
	}

	SendResponse(ctx, 200, nil, "Added members to group", nil)
}

func GetUsersOfGroups(ctx *gin.Context) {
	_, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, 401, nil, "Unauthorized", errors.New("unauthorized"))
		return
	}

	// TODO: Check if user is authorized to see members of the group
	var req dto.GetUsersOfGroupsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	groupUsers, err := service.GetUsersOfGroups(ctx, req.GroupIDs)

	if err != nil {
		SendResponse(ctx, 500, nil, "failed to fetch group members", nil)
		return
	}
	SendResponse(ctx, 200, groupUsers, "Fetched group users", nil)
}

func GetCredentialGroups(ctx *gin.Context) {
	credentialIDStr := ctx.Param("id")
	credentialID, _ := uuid.Parse(credentialIDStr)

	users, err := service.GetCredentialGroups(ctx, credentialID)
	if err != nil {
		SendResponse(ctx, 400, nil, "failed to get credential users", err)
		return
	}
	SendResponse(ctx, 200, users, "fetched credential users", nil)
}

func GetUsersWithoutGroupAccess(ctx *gin.Context) {
	userID, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, 401, nil, "Unauthorized", errors.New("unauthorized"))
		return
	}

	groupIDStr := ctx.Param("groupId")
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		logger.Errorf(err.Error())
		SendResponse(ctx, 400, nil, "Invalid group id", errors.New("invalid group id"))
		return
	}

	users, err := service.GetUsersWithoutGroupAccess(ctx, userID, groupID)

	if err != nil {
		logger.Errorf(err.Error())
		SendResponse(ctx, 500, nil, "", err)
		return
	}
	SendResponse(ctx, 200, users, "Fetched users not in group", nil)
}
