package controllers

import (
	"errors"
	"net/http"
	dto "osvauld/dtos"
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
		SendResponse(ctx, 400, nil, "failed to create user", errors.New("failed to create user"))
		return
	}
	SendResponse(ctx, 201, user, "created user", nil)

}

func Login(ctx *gin.Context) {
	var req dto.Login

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := service.Login(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	SendResponse(ctx, 200, user, "Login successfull", nil)
}

func GetAllUsers(ctx *gin.Context) {
	users, _ := service.GetAllUsers(ctx)
	SendResponse(ctx, 200, users, "fetched users", nil)
}

func GetChallenge(ctx *gin.Context) {
	var req dto.CreateChallenge
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	challenge, _ := service.CreateChallenge(ctx, req)
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
	token, _ := service.VerifyChallenge(ctx, req)
	type TokenResponse struct {
		Token string `json:"token"`
	}
	SendResponse(ctx, 200, TokenResponse{Token: token}, "verified challenge", nil)
}
