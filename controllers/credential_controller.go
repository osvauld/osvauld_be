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
		SendResponse(ctx, http.StatusBadRequest, nil, "", err)
		return
	}

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}

	credentialID, err := service.AddCredential(ctx, req, caller)
	if err != nil {
		logger.Errorf(err.Error())

		if _, ok := err.(*customerrors.UserNotManagerOfFolderError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", err)
			return
		}
		SendResponse(ctx, http.StatusInternalServerError, nil, "Failed to add credential", nil)
		return
	}

	response := map[string]uuid.UUID{"credentialId": credentialID}
	SendResponse(ctx, http.StatusOK, response, "Added credential", nil)
}

func GetCredentialDataByID(ctx *gin.Context) {
	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}

	credentialIDStr := ctx.Param("id")
	credentailaID, err := uuid.Parse(credentialIDStr)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid credential id"))
		return
	}

	credential, err := service.GetCredentialDataByID(ctx, credentailaID, caller)
	if err != nil {
		if _, ok := err.(*customerrors.UserDoesNotHaveCredentialAccessError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", errors.New("unauthorized"))
			return
		}
		SendResponse(ctx, http.StatusInternalServerError, nil, "Failed to fetch credential", errors.New("failed to fetch credential"))
		return
	}
	SendResponse(ctx, http.StatusOK, credential, "Fetched credential", nil)
}

func GetCredentialsByFolder(ctx *gin.Context) {
	// Parse user_id from headerss

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}
	// Get folder_id from query params
	folderIDStr := ctx.Param("id")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid folder id"))
		return
	}

	credentials, err := service.GetCredentialsByFolder(ctx, folderID, caller)
	if err != nil {
		SendResponse(ctx, 500, nil, "", errors.New("failed to fetch credential"))
		return
	}
	SendResponse(ctx, 200, credentials, "Fetched credentials", nil)

}

func GetCredentialsFieldsByFolderID(ctx *gin.Context) {
	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}

	folderIDStr := ctx.Param("folderId")
	folderID, err := uuid.Parse(folderIDStr)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid folder id"))
		return
	}

	credential, err := service.GetCredentialsFieldsByFolderID(ctx, folderID, caller)
	if err != nil {
		if _, ok := err.(*customerrors.UserDoesNotHaveFolderAccessError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", errors.New("unauthorized"))
			return
		}
		SendResponse(ctx, http.StatusInternalServerError, nil, "", errors.New("failed to fetch credential"))
		return
	}
	SendResponse(ctx, http.StatusOK, credential, "Fetched credential", nil)
}

func GetCredentialsFieldsByIds(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}

	var req dto.GetCredentialsFieldsByIdsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	credentials, err := service.GetFieldsByCredentialIDs(ctx, req.CredentialIds, caller)
	if err != nil {
		if _, ok := err.(*customerrors.UserDoesNotHaveCredentialAccessError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", errors.New("unauthorized"))
			return
		}

		SendResponse(ctx, http.StatusInternalServerError, nil, "Failed to fetch credential", errors.New("failed to fetch credential"))
		return
	}
	SendResponse(ctx, http.StatusOK, credentials, "Fetched credential", nil)
}

// This is used by search to get credentials by ids
func GetCredentialsByIDs(ctx *gin.Context) {
	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}

	var req dto.GetCredentialsByIDsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", err)
		return
	}

	credentials, err := service.GetCredentialsByIDs(ctx, req.CredentialIds, caller)
	if err != nil {

		if _, ok := err.(*customerrors.UserDoesNotHaveCredentialAccessError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", errors.New("unauthorized"))
			return
		}

		SendResponse(ctx, http.StatusInternalServerError, nil, "Failed to fetch credential", errors.New("failed to fetch credential"))
		return
	}
	SendResponse(ctx, http.StatusOK, credentials, "Fetched credential", nil)
}

func GetAllUrlsForUser(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusUnauthorized, nil, "", errors.New("unauthorized"))
		return
	}

	urls, err := service.GetAllUrlsForUser(ctx, caller)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, nil, "", errors.New("failed to fetch urls"))
		return
	}
	SendResponse(ctx, http.StatusOK, urls, "Fetched urls", nil)
}

func GetSensitiveFieldsByCredentialID(ctx *gin.Context) {

	userID, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusUnauthorized, nil, "", errors.New("unauthorized"))
		return
	}

	credentialIDStr := ctx.Param("id")
	credentailID, err := uuid.Parse(credentialIDStr)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid credential id"))
		return
	}

	sensitiveFields, err := service.GetSensitiveFields(ctx, credentailID, userID)
	if err != nil {

		if _, ok := err.(*customerrors.UserDoesNotHaveCredentialAccessError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", err)
			return
		}

		SendResponse(ctx, http.StatusInternalServerError, nil, "Failed to fetch credential", errors.New("failed to fetch credential"))
		return
	}

	SendResponse(ctx, http.StatusOK, sensitiveFields, "Fetched credential", nil)

}

func EditCredential(ctx *gin.Context) {
	var req dto.EditCredentialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", err)
		return
	}

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}

	credentialIDStr := ctx.Param("id")
	credentailID, err := uuid.Parse(credentialIDStr)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid credential id"))
		return
	}

	err = service.EditCredential(ctx, credentailID, req, caller)
	if err != nil {
		logger.Errorf(err.Error())

		if _, ok := err.(*customerrors.UserNotManagerOfCredentialError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", err)
			return
		}
		SendResponse(ctx, http.StatusInternalServerError, nil, "Failed to edit credential", nil)
		return
	}

	SendResponse(ctx, http.StatusOK, nil, "", nil)
}

func EditCredentialDetails(ctx *gin.Context) {
	var req dto.EditCredentialDetailsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", err)
		return
	}

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}

	credentialIDStr := ctx.Param("id")
	credentailID, err := uuid.Parse(credentialIDStr)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid credential id"))
		return
	}

	err = service.EditCredentialDetails(ctx, credentailID, req, caller)
	if err != nil {
		logger.Errorf(err.Error())

		if _, ok := err.(*customerrors.UserNotManagerOfCredentialError); ok {
			SendResponse(ctx, http.StatusUnauthorized, nil, "", err)
			return
		}
		SendResponse(ctx, http.StatusInternalServerError, nil, "Failed to edit credential", nil)
		return
	}

	SendResponse(ctx, http.StatusOK, nil, "", nil)
}

func GetSearchData(ctx *gin.Context) {

	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}
	credentials, err := service.GetSearchData(ctx, caller)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}
	SendResponse(ctx, http.StatusOK, credentials, "Fetched credential", nil)
}

func RemoveCredential(ctx *gin.Context) {
	caller, err := utils.FetchUserIDFromCtx(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid user id"))
		return
	}
	credentialIDStr := ctx.Param("id")
	credentailID, err := uuid.Parse(credentialIDStr)
	if err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", errors.New("invalid credential id"))
		return
	}

	err = service.RemoveCredential(ctx, credentailID, caller)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, nil, "", errors.New("failed to remove credential"))
		return
	}
	SendResponse(ctx, http.StatusOK, nil, "Credential removed successfully", nil)
}
