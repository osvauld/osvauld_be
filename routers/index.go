package routers

import (
	"net/http"
	"osvauld/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes add all routing list here automatically get main router
func RegisterRoutes(route *gin.Engine) {
	route.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Route Not Found"})
	})
	route.GET("/health", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"live": "ok"}) })

	route.POST("/user/", controllers.CreateUser)
	route.POST("/folder/", controllers.CreateFolder)
	route.POST("/credential/", controllers.AddCredential)
	route.GET("/secrets/", controllers.GetCredentialsByFolder)
	route.PUT("/secrets/", controllers.ShareCredential)
	route.GET("/folders/", controllers.GetAccessibleFolders)
	// route.GET("/credential/:id", controllers.GetCredentialByID)
	route.POST("/group", controllers.AddGroup)
	route.POST("/group/members", controllers.AppendMembersToGroup)
	//Add All route
	//TestRoutes(route)
}
