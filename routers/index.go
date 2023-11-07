package routers

import (
	"net/http"
	"osvauld/controllers"
	"osvauld/routers/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes add all routing list here automatically get main router
func RegisterRoutes(route *gin.Engine) {
	route.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Route Not Found"})
	})
	route.GET("/health", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"live": "ok"}) })
	route.POST("/login", controllers.Login)
	route.POST("/user/", middleware.JWTAuthMiddleware(), controllers.CreateUser)
	route.POST("/folder/", middleware.JWTAuthMiddleware(), controllers.CreateFolder)
	route.POST("/credential/", middleware.JWTAuthMiddleware(), controllers.AddCredential)
	route.GET("/secrets/", middleware.JWTAuthMiddleware(), controllers.GetCredentialsByFolder)
	route.PUT("/secrets/", middleware.JWTAuthMiddleware(), controllers.ShareCredential)
	route.GET("/folders/", middleware.JWTAuthMiddleware(), controllers.GetAccessibleFolders)
	route.GET("/credential/:id", middleware.JWTAuthMiddleware(), controllers.GetCredentialByID)
	route.POST("/group", middleware.JWTAuthMiddleware(), controllers.AddGroup)
	route.GET("/group/:groupId", middleware.JWTAuthMiddleware(), controllers.GetGroupMembers)
	route.GET("/groups", middleware.JWTAuthMiddleware(), controllers.GetUserGroups)
	route.POST("/group/members", middleware.JWTAuthMiddleware(), controllers.AppendMembersToGroup)
	//Add All route
	//TestRoutes(route)
}
