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
	"github.com/google/uuid"
)

func AddCredential(ctx *gin.Context) {

	var req dto.AddCredentialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, 401, nil, "Unauthorized", errors.New("unauthorized"))
		return
	}

	credentialID, err := service.AddCredential(ctx, req, caller)
	if err != nil {
		logger.Errorf(err.Error())

		if _, ok := err.(*customerrors.UserNotAuthenticatedError); ok {
			SendResponse(ctx, 401, nil, "", err)
			return
		}
		SendResponse(ctx, 500, nil, "Failed to add credential", nil)
		return
	}

	response := map[string]uuid.UUID{"credentialId": credentialID}
	SendResponse(ctx, 200, response, "Added credential", nil)
}

func GetCredentialDataByID(ctx *gin.Context) {
	userIdInterface, _ := ctx.Get("userId")
	userID, _ := userIdInterface.(uuid.UUID)
	credentialIDStr := ctx.Param("id")
	credentailaID, _ := uuid.Parse(credentialIDStr)
	credential, err := service.GetCredentialDataByID(ctx, credentailaID, userID)
	if err != nil {
		if _, ok := err.(*customerrors.UserNotAuthenticatedError); ok {
			SendResponse(ctx, 401, nil, "Unauthorized", errors.New("unauthorized"))
			return
		}
		SendResponse(ctx, 200, nil, "Failed to fetch credential", errors.New("failed to fetch credential"))
		return
	}
	SendResponse(ctx, 200, credential, "Fetched credential", nil)
}

func GetCredentialsByFolder(ctx *gin.Context) {
	// Parse user_id from headerss

	userID, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, 401, nil, "Unauthorized", errors.New("unauthorized"))
		return
	}
	// Get folder_id from query params
	folderIDStr := ctx.Param("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		SendResponse(ctx, 400, nil, "Invalid folder id", errors.New("invalid folder id"))
		return
	}

	credentials, err := service.GetCredentialsByFolder(ctx, folderID, userID)
	if err != nil {
		SendResponse(ctx, 500, nil, "Failed to fetch credential", errors.New("failed to fetch credential"))
		return
	}
	SendResponse(ctx, 200, credentials, "Fetched credentials", nil)

}

func GetCredentialsFieldsByFolderID(ctx *gin.Context) {
	userIdInterface, _ := ctx.Get("userId")
	userID, _ := userIdInterface.(uuid.UUID)
	folderIDStr := ctx.Param("folderId")
	folderID, _ := uuid.Parse(folderIDStr)
	credential, err := service.GetCredentialsFieldsByFolderID(ctx, folderID, userID)
	if err != nil {
		SendResponse(ctx, 200, nil, "Failed to fetch credential", errors.New("failed to fetch credential"))
		return
	}
	SendResponse(ctx, 200, credential, "Fetched credential", nil)
}

func GetCredentialsFieldsByIds(ctx *gin.Context) {
	userIdInterface, _ := ctx.Get("userId")
	userID, _ := userIdInterface.(uuid.UUID)
	var req dto.GetCredentialsFieldsByIdsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	credentials, err := service.GetFieldsByCredentialIds(ctx, req.CredentialIds, userID)
	if err != nil {
		SendResponse(ctx, 200, nil, "Failed to fetch credential", errors.New("failed to fetch credential"))
		return
	}
	SendResponse(ctx, 200, credentials, "Fetched credential", nil)
}

// This is used by search to get credentials by ids
func GetCredentialsByIDs(ctx *gin.Context) {
	userIdInterface, _ := ctx.Get("userId")
	userID, _ := userIdInterface.(uuid.UUID)
	var req dto.GetCredentialsByIDsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	credentials, err := service.GetCredentialsByIDs(ctx, req.CredentialIds, userID)
	if err != nil {
		SendResponse(ctx, 200, nil, "Failed to fetch credential", errors.New("failed to fetch credential"))
		return
	}
	SendResponse(ctx, 200, credentials, "Fetched credential", nil)
}

func GetAllUrlsForUser(ctx *gin.Context) {
	userIdInterface, _ := ctx.Get("userId")
	userID, _ := userIdInterface.(uuid.UUID)
	urls, err := service.GetAllUrlsForUser(ctx, userID)
	if err != nil {
		SendResponse(ctx, 200, nil, "Failed to fetch urls", errors.New("failed to fetch urls"))
		return
	}
	SendResponse(ctx, 200, urls, "Fetched urls", nil)
}

func GetSensitiveFieldsCredentialByID(ctx *gin.Context) {

	userID, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, 401, nil, "", errors.New("unauthorized"))
		return
	}

	credentialIDStr := ctx.Param("id")
	credentailID, err := uuid.Parse(credentialIDStr)
	if err != nil {
		SendResponse(ctx, 400, nil, "", errors.New("invalid credential id"))
		return
	}

	sensitiveFields, err := service.GetSensitiveFields(ctx, credentailID, userID)
	if err != nil {

		if _, ok := err.(*customerrors.UserNotAuthenticatedError); ok {
			SendResponse(ctx, 401, nil, "", err)
			return
		}

		SendResponse(ctx, 500, nil, "Failed to fetch credential", errors.New("failed to fetch credential"))
		return
	}

	SendResponse(ctx, 200, sensitiveFields, "Fetched credential", nil)

}

func EditCredential(ctx *gin.Context) {
	var req dto.EditCredentialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, 401, nil, "Unauthorized", errors.New("unauthorized"))
		return
	}

	credentialIDStr := ctx.Param("id")
	credentailID, err := uuid.Parse(credentialIDStr)
	if err != nil {
		SendResponse(ctx, 400, nil, "Invalid credential id", errors.New("invalid credential id"))
		return
	}

	err = service.EditCredential(ctx, credentailID, req, caller)
	if err != nil {
		logger.Errorf(err.Error())

		if _, ok := err.(*customerrors.UserNotAuthenticatedError); ok {
			SendResponse(ctx, 401, nil, "", err)
			return
		}
		SendResponse(ctx, 500, nil, "Failed to edit credential", nil)
		return
	}

	SendResponse(ctx, 200, nil, "", nil)
}
