package controllers

import (
	"net/http"
	dto "osvauld/dtos"
	"osvauld/infra/logger"
	service "osvauld/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

// // Initialize the validator once for your application
var validate = validator.New()

func CreateUser(ctx *gin.Context) {
	// Define a struct just for request body

	//TODO: add created by field
	var req dto.CreateUser
	// Bind the request body to the struct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", err)
		return
	}

	// Validate the requestBody using the validator
	if err := validate.Struct(req); err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", err)
		return
	}
	user, err := service.CreateUser(ctx, req)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}
	SendResponse(ctx, http.StatusCreated, user, "created user", nil)

}

func TempLogin(ctx *gin.Context) {

	var req dto.TempLogin
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", err)
		return
	}
	challenge, err := service.TempLogin(ctx, req)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}

	response := map[string]string{"challenge": challenge}

	SendResponse(ctx, http.StatusOK, response, "temp login successfull", nil)
}

// single use api per user to register their public keys
func Register(ctx *gin.Context) {
	var req dto.Register
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "invalid request", err)
		return
	}
	user, err := service.Register(ctx, req)
	if err != nil {
		logger.Errorf(err.Error())
		SendResponse(ctx, http.StatusInternalServerError, nil, "failed to register user", err)
		return
	}

	SendResponse(ctx, http.StatusOK, user, "registration successfull", nil)
}

func GetChallenge(ctx *gin.Context) {
	var req dto.CreateChallenge
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", err)
		return
	}
	challenge, err := service.CreateChallenge(ctx, req)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}
	type ChallengeResponse struct {
		Challenge string `json:"challenge"`
	}
	SendResponse(ctx, http.StatusOK, ChallengeResponse{Challenge: challenge}, "fetched challenge", nil)
}

func VerifyChallenge(ctx *gin.Context) {
	var req dto.VerifyChallenge
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", err)
		return
	}
	token, err := service.VerifyChallenge(ctx, req)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}
	type TokenResponse struct {
		Token string `json:"token"`
	}
	SendResponse(ctx, http.StatusOK, TokenResponse{Token: token}, "verified challenge", nil)
}

func GetAllUsers(ctx *gin.Context) {
	users, err := service.GetAllUsers(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}
	SendResponse(ctx, http.StatusOK, users, "fetched users", nil)
}

func GetAdminPage(ctx *gin.Context) {
	exists, err := service.CheckUserExists(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}
	if exists {
		ctx.HTML(http.StatusOK, "admin_exists.tmpl", nil)
	} else {
		ctx.HTML(http.StatusOK, "admin_create.tmpl", nil)
	}
}

func CreateFirstAdmin(ctx *gin.Context) {
	exists, err := service.CheckUserExists(ctx)
	if err != nil {
		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}

	if exists {
		// Admin user already exists, render the "user exists" template
		ctx.HTML(http.StatusOK, "admin_exists.tmpl", nil)
		return
	}

	var req dto.CreateUser

	// Bind the request body to the struct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", err)
		return
	}

	// Validate the requestBody using the validator
	if err := validate.Struct(req); err != nil {
		SendResponse(ctx, http.StatusBadRequest, nil, "", err)
		return
	}

	_, err = service.CreateUser(ctx, req)
	if err != nil {
		logger.Errorf(err.Error())
		SendResponse(ctx, http.StatusInternalServerError, nil, "", err)
		return
	}

	// Admin user created successfully, render the "admin created" template
	ctx.HTML(http.StatusOK, "admin_created.tmpl", nil)
}

func RemoveUserFromAll(ctx *gin.Context) {
	// TODO: check user is deleteing themselves
	id := ctx.Param("id")
	userID, _ := uuid.Parse(id)
	err := service.RemoveUserFromAll(ctx, userID)
	if err != nil {
		SendResponse(ctx, 400, nil, "failed to delete user", err)
		return
	}
	SendResponse(ctx, 200, nil, "deleted user", nil)
}
