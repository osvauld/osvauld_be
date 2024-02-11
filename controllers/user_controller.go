package controllers

import (
	"errors"
	"net/http"
	dto "osvauld/dtos"
	"osvauld/infra/logger"
	service "osvauld/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

// // Initialize the validator once for your application
var validate = validator.New()

func CreateUser(ctx *gin.Context) {
	// Define a struct just for request body

	//TODO: add created by field
	var req dto.CreateUser
	// Bind the request body to the struct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the requestBody using the validator
	if err := validate.Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := service.CreateUser(ctx, req)
	if err != nil {
		logger.Errorf(err.Error())
		SendResponse(ctx, 400, nil, "failed to create user", nil)
		return
	}
	SendResponse(ctx, 201, user, "created user", nil)

}

// single use api per user to register their public keys
func Register(ctx *gin.Context) {
	var req dto.Register
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := service.Register(ctx, req)
	if err != nil {
		SendResponse(ctx, 400, nil, "failed to register user", errors.New("failed to register user"))
		return
	}

	SendResponse(ctx, 200, user, "registration successfull", nil)
}

func GetChallenge(ctx *gin.Context) {
	var req dto.CreateChallenge
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	challenge, err := service.CreateChallenge(ctx, req)
	if err != nil {
		SendResponse(ctx, 400, nil, "", err)
		return
	}
	type ChallengeResponse struct {
		Challenge string `json:"challenge"`
	}
	SendResponse(ctx, 200, ChallengeResponse{Challenge: challenge}, "fetched challenge", nil)
}

func VerifyChallenge(ctx *gin.Context) {
	var req dto.VerifyChallenge
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := service.VerifyChallenge(ctx, req)
	if err != nil {
		SendResponse(ctx, 400, nil, "", err)
		return
	}
	type TokenResponse struct {
		Token string `json:"token"`
	}
	SendResponse(ctx, 200, TokenResponse{Token: token}, "verified challenge", nil)
}

func GetAllUsers(ctx *gin.Context) {
	users, _ := service.GetAllUsers(ctx)
	SendResponse(ctx, 200, users, "fetched users", nil)
}

func GetCredentialUsers(ctx *gin.Context) {
	credentialIDStr := ctx.Param("id")
	credentialID, _ := uuid.Parse(credentialIDStr)

	users, err := service.GetCredentialUsers(ctx, credentialID)
	if err != nil {
		SendResponse(ctx, 400, nil, "failed to get credential users", err)
		return
	}
	SendResponse(ctx, 200, users, "fetched credential users", nil)
}
