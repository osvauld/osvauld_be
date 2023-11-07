package controllers

import (
	"errors"
	"net/http"
	dto "osvauld/dtos"
	service "osvauld/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddCredential(ctx *gin.Context) {
	var req dto.AddCredentailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIdInterface, _ := ctx.Get("userId")
	userID, _ := userIdInterface.(uuid.UUID)
	credentialId, err := service.CreateCredential(ctx, req, userID)
	if err != nil {
		SendResponse(ctx, 500, nil, "Failed to add credential", errors.New("failed to add credential"))
	}
	SendResponse(ctx, 201, credentialId, "Added Credential", nil)
}

func GetCredentialsByFolder(ctx *gin.Context) {
	// Parse user_id from headerss
	userIdInterface, _ := ctx.Get("userId")
	userID, _ := userIdInterface.(uuid.UUID)
	// Get folder_id from query params
	folderIDStr := ctx.DefaultQuery("folderId", "")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder_id query parameter."})
		return
	}
	credentials, err := service.GetCredentialsByFolder(ctx, folderID, userID)
	if err != nil {
		SendResponse(ctx, 500, nil, "Failed to fetch credential", errors.New("Failed to fetch credential"))
		return
	}
	SendResponse(ctx, 200, credentials, "Fetched credentials", nil)

}

func ShareCredential(ctx *gin.Context) {
	var req dto.ShareCredentialPayload
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIdInterface, _ := ctx.Get("userId")
	userID, _ := userIdInterface.(uuid.UUID)
	service.ShareCredential(ctx, req, userID)
	SendResponse(ctx, 200, nil, "Success", nil)
}

func GetCredentialByID(ctx *gin.Context) {
	userIdInterface, _ := ctx.Get("userId")
	userID, _ := userIdInterface.(uuid.UUID)
	credentialIDStr := ctx.Param("id")
	credentailaID, _ := uuid.Parse(credentialIDStr)
	credential, err := service.FetchCredentialByID(ctx, credentailaID, userID)
	if err != nil {
		SendResponse(ctx, 200, nil, "Failed to fetch credential", errors.New("failed to fetch credential"))
		return
	}
	SendResponse(ctx, 200, credential, "Fetched credential", nil)
}
