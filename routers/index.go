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
	route.POST("/user/", controllers.CreateUser)
	route.GET("/users", middleware.JWTAuthMiddleware(), controllers.GetAllUsers)
	route.POST("/folder/", middleware.JWTAuthMiddleware(), controllers.CreateFolder)
	route.PUT("/folder", middleware.JWTAuthMiddleware(), controllers.ShareFolder)
	route.GET("/folder/:id", middleware.JWTAuthMiddleware(), controllers.GetUsersByFolder)
	route.GET("/folder/:id/users", middleware.JWTAuthMiddleware(), controllers.GetSharedUsers)
	route.GET("/folder/:id/credential", middleware.JWTAuthMiddleware(), controllers.GetCredentialsByFolder)
	route.PUT("/credentials/", middleware.JWTAuthMiddleware(), controllers.ShareCredential)
	route.GET("/folders/", middleware.JWTAuthMiddleware(), controllers.GetAccessibleFolders)

	// Credential Routes
	route.POST("/credential/", middleware.JWTAuthMiddleware(), controllers.AddCredential)
	route.GET("/credential/:id", middleware.JWTAuthMiddleware(), controllers.GetCredentialByID)

	// route.GET("/credential/:id", middleware.JWTAuthMiddleware(), controllers.GetCredentialByID)
	route.POST("/group", middleware.JWTAuthMiddleware(), controllers.AddGroup)
	route.GET("/group/:groupId", middleware.JWTAuthMiddleware(), controllers.GetGroupMembers)
	route.GET("/groups", middleware.JWTAuthMiddleware(), controllers.GetUserGroups)
	route.POST("/group/members", middleware.JWTAuthMiddleware(), controllers.AppendMembersToGroup)
	route.GET("/credentials/encrypted/:folderId", middleware.JWTAuthMiddleware(), controllers.GetEncryptedCredentails)
	route.POST("/credentials/encrypted/", middleware.JWTAuthMiddleware(), controllers.GetEncryptedCredentailsByIds)
	//Add All route
	//TestRoutes(route)
}
