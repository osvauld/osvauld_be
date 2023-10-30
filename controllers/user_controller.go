package controllers

import (
	"fmt"
	"net/http"
	"osvauld/models"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Initialize the validator once for your application
var validate = validator.New()

func CreateUser(ctx *gin.Context) {
	// Define a struct just for request body
	var requestBody struct {
		Username string `json:"username" binding:"required,min=3,max=32"`
	}

	// Bind the request body to the struct
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the requestBody using the validator
	if err := validate.Struct(requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a User object
	user := &models.User{
		Username: requestBody.Username,
	}
	fmt.Println(user)
	// Save the user to the database
	if err := repository.SaveUser(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user."})
		return
	}
	ctx.JSON(http.StatusOK, user)
}
