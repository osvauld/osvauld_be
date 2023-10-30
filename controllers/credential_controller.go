package controllers

import (
	"net/http"
	"osvauld/infra/logger"
	"osvauld/models"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SecretRequest struct {
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	FolderID          uuid.UUID      `json:"folder_id"`
	EncryptedFields   []FieldRequest `json:"encrypted_fields"`
	UnencryptedFields []FieldRequest `json:"unencrypted_fields"`
}

type FieldRequest struct {
	FieldName  string `json:"field_name"`
	FieldValue string `json:"field_value"`
}

func AddSecret(ctx *gin.Context) {
	var req SecretRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user_id := ctx.GetHeader("user_id")
	if user_id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id header is required"})
		return
	}
	user_uuid, _ := uuid.Parse(user_id)
	// Create the Credential
	credential := models.Credential{
		Name:        req.Name,
		Description: req.Description,
		FolderID:    req.FolderID,
		CreatedBy:   user_uuid,
	}

	err := repository.SaveCredential(&credential)
	if err != nil {
		logger.Errorf(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save credential."})
		return
	}
	access_list := models.AccessList{
		CredentialID: credential.ID,
		UserID:       user_uuid,
		AccessType:   "OWNER",
	}
	err = repository.AddAccessList(&access_list)
	if err != nil {
		logger.Errorf(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save credential."})
		return
	}
	// Save Encrypted Fields
	for _, field := range req.EncryptedFields {
		encryptedData := models.EncryptedData{
			FieldName:    field.FieldName,
			CredentialID: credential.ID,
			FieldValue:   field.FieldValue,
			UserID:       user_uuid,
		}

		err = repository.SaveEncryptedData(&encryptedData)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save encrypted data."})
			return
		}
	}

	// Save Unencrypted Fields
	for _, field := range req.UnencryptedFields {
		unencryptedData := models.UnencryptedData{
			FieldName:    field.FieldName,
			CredentialID: credential.ID,
			FieldValue:   field.FieldValue,
		}

		err = repository.SaveUnencryptedData(&unencryptedData)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save unencrypted data."})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Secret successfully saved!"})
}
