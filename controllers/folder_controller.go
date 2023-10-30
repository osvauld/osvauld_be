package controllers

import (
	"net/http"
	"osvauld/models"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateFolder(ctx *gin.Context) {
	var folder models.Folder
	if err := ctx.ShouldBindJSON(&folder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validation using go-playground/validator
	if err := validate.Struct(folder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve user_id from header
	user_id := ctx.GetHeader("user_id")
	if user_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id header is required"})
		return
	}
	folder.CreatedBy, _ = uuid.Parse(user_id)
	if err := repository.SaveFolder(&folder); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create folder."})
		return
	}
	ctx.JSON(http.StatusOK, folder)
}
