package controllers

import (
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
	user, _ := service.CreateUser(ctx, req)

	ctx.JSON(http.StatusOK, user)
}
