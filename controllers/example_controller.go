package controllers

import (
	"net/http"
	"osvauld/models"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
)

func GetData(ctx *gin.Context) {
	var example []*models.Example
	repository.Get(&example)
	ctx.JSON(http.StatusOK, &example)

}
