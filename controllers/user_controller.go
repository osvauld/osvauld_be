package controllers

import (
	"errors"
	"net/http"
	dto "osvauld/dtos"
	service "osvauld/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

// // Initialize the validator once for your application
var validate = validator.New()

func CreateUser(ctx *gin.Context) {
	// Define a struct just for request body

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
	user := service.Login(ctx, req)
	SendResponse(ctx, 200, user, "Login successfull", nil)
}
