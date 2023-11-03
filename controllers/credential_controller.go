package controllers

import (
	"net/http"
	dto "osvauld/dtos"
	"osvauld/infra/logger"
	service "osvauld/services"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddCredential(ctx *gin.Context) {
	var req dto.AddCredentailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIdString := ctx.GetHeader("userId")
	if userIdString == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id header is required"})
		return
	}
	userID, err := uuid.Parse(strings.Trim(userIdString, `"`))
	if err != nil {
		logger.Errorf(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id header is required"})
		return
	}
	service.CreateCredential(ctx, req, userID)
	ctx.JSON(http.StatusOK, gin.H{"message": "Secret successfully saved!"})
}

func GetCredentialsByFolder(ctx *gin.Context) {
	// Parse user_id from header
	userIDHeader := ctx.GetHeader("userId")
	userID, err := uuid.Parse(userIDHeader)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id in header."})
		return
	}

	// Get folder_id from query params
	folderIDStr := ctx.DefaultQuery("folderId", "")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder_id query parameter."})
		return
	}
	credentials, _ := service.GetCredentialsByFolder(ctx, folderID, userID)

	ctx.JSON(http.StatusOK, credentials)
}

func ShareCredential(ctx *gin.Context) {
	var req dto.ShareCredentialPayload
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDHeader := ctx.GetHeader("userId")
	userID, err := uuid.Parse(userIDHeader)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id in header."})
		return
	}
	service.ShareCredential(ctx, req, userID)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func GetCredentialByID(ctx *gin.Context) {

	userIDHeader := ctx.GetHeader("userId")
	userID, err := uuid.Parse(userIDHeader)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id in header."})
		return
	}
	credentialIDStr := ctx.Param("id")
	credentailaID, _ := uuid.Parse(credentialIDStr)
	credential, _ := service.FetchCredentialByID(ctx, credentailaID, userID)
	ctx.JSON(http.StatusOK, credential)
}
