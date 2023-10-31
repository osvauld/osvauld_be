package routers

import (
	"net/http"
	"osvauld/controllers"

	"github.com/gin-gonic/gin"
)

//RegisterRoutes add all routing list here automatically get main router
func RegisterRoutes(route *gin.Engine) {
	route.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Route Not Found"})
	})
	route.GET("/health", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"live": "ok"}) })

	route.POST("/v1/user/", controllers.CreateUser)
	route.POST("/v1/folder/", controllers.CreateFolder)
	route.POST("/v1/credential/", controllers.AddSecret)
	route.GET("/secrets/", controllers.GetSecretsForUser)
	route.PUT("/secrets/", controllers.ShareSecret)
	route.GET("/folders/", controllers.GetAccessibleFolders)
	//Add All route
	//TestRoutes(route)
}
