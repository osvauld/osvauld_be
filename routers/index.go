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
	// route.POST("/login", controllers.Login)
	route.POST("/user/", controllers.CreateUser)
	route.POST("/user/register", controllers.Register)
	route.POST("/user/challenge", controllers.GetChallenge)
	route.POST("/user/verify", controllers.VerifyChallenge)
	route.GET("/users", middleware.JWTAuthMiddleware(), controllers.GetAllUsers)
	route.POST("/folder/", middleware.JWTAuthMiddleware(), controllers.CreateFolder)
	route.PUT("/folder", middleware.JWTAuthMiddleware(), controllers.ShareFolder)
	route.GET("/folder/:id", middleware.JWTAuthMiddleware(), controllers.GetUsersByFolder)
	route.GET("/folder/:id/users", middleware.JWTAuthMiddleware(), controllers.GetSharedUsers)
	route.GET("/folder/:id/credential", middleware.JWTAuthMiddleware(), controllers.GetCredentialsByFolder)
	route.GET("/folders/", middleware.JWTAuthMiddleware(), controllers.GetAccessibleFolders)

	route.POST("/shareCredential/Users", middleware.JWTAuthMiddleware(), controllers.ShareMultipleCredentialsWithMulitpleUsers)

	// Credential Routes
	route.POST("/credential/", middleware.JWTAuthMiddleware(), controllers.AddCredential)
	route.GET("/credential/:id", middleware.JWTAuthMiddleware(), controllers.FetchCredentialByID)

	// route.GET("/credential/:id", middleware.JWTAuthMiddleware(), controllers.GetCredentialByID)
	route.POST("/group", middleware.JWTAuthMiddleware(), controllers.CreateGroup)
	route.GET("/group/:groupId", middleware.JWTAuthMiddleware(), controllers.GetGroupMembers)
	route.GET("/groups", middleware.JWTAuthMiddleware(), controllers.GetUserGroups)
	route.POST("/group/members", middleware.JWTAuthMiddleware(), controllers.AddMemberToGroup)

	route.GET("/credentials/encrypted/:folderId", middleware.JWTAuthMiddleware(), controllers.GetAllEncryptedCredentailsForFolderID)
	route.POST("/credentials/encrypted/", middleware.JWTAuthMiddleware(), controllers.GetEncryptedCredentailsByIds)

	route.GET("/group/:groupId/encrypted/", middleware.JWTAuthMiddleware(), controllers.FetchEncryptedValuesWithGroupAccess)

	//Add All route
	//TestRoutes(route)
}
