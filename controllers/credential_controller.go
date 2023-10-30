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

type ShareSecretPayload struct {
	CredentialList []struct {
		CredentialID string `json:"credential_id"`
		Users        []struct {
			UserID string `json:"user_id"`
			Fields []struct {
				FieldName  string `json:"field_name"`
				FieldValue string `json:"field_value"`
			} `json:"fields"`
			AccessType string `json:"access_type"`
		} `json:"users"`
	} `json:"credential_list"`
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

func GetSecretsForUser(ctx *gin.Context) {
	// Parse user_id from header
	userIDHeader := ctx.GetHeader("user_id")
	userID, err := uuid.Parse(userIDHeader)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id in header."})
		return
	}

	// Get folder_id from query params
	folderIDStr := ctx.DefaultQuery("folder_id", "")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder_id query parameter."})
		return
	}

	// Fetch secrets
	secrets, err := repository.GetSecretsByFolderAndUser(folderID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve secrets."})
		return
	}

	var response []map[string]interface{}

	for _, secret := range secrets {
		data := make([]map[string]string, 0)
		for _, field := range secret.Fields {
			data = append(data, map[string]string{
				"FieldName":  field.FieldName,
				"FieldValue": field.FieldValue,
			})
		}

		response = append(response, map[string]interface{}{
			"ID":     secret.ID,
			"Fields": data,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

func ShareSecret(ctx *gin.Context) {
	var payload ShareSecretPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Iterate over each credential in the payload
	for _, credential := range payload.CredentialList {
		credential_id, _ := uuid.Parse(credential.CredentialID)

		// Iterate over users for this credential
		for _, user := range credential.Users {
			user_id, _ := uuid.Parse(user.UserID)

			// Insert into AccessList
			access_list := models.AccessList{
				CredentialID: credential_id,
				UserID:       user_id,
				AccessType:   user.AccessType,
			}
			err := repository.AddAccessList(&access_list)
			if err != nil {
				logger.Errorf(err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save credential."})
				return
			}

			// Insert fields into EncryptedData table
			for _, field := range user.Fields {
				encryptedData := models.EncryptedData{
					FieldName:    field.FieldName,
					CredentialID: credential_id,
					FieldValue:   field.FieldValue,
					UserID:       user_id,
				}
				err := repository.SaveEncryptedData(&encryptedData)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save encrypted data."})
					return
				}
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
