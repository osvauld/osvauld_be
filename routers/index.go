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
	route.POST("/testLogin", controllers.TestLoginController)
	route.POST("/user/", controllers.CreateUser)
	route.POST("/user/register", controllers.Register)
	route.POST("/user/challenge", controllers.GetChallenge)
	route.POST("/user/verify", controllers.VerifyChallenge)
	route.GET("/users", middleware.JWTAuthMiddleware(), controllers.GetAllUsers)
	route.POST("/folder/", middleware.JWTAuthMiddleware(), controllers.CreateFolder)
	route.GET("/folder/:id", middleware.JWTAuthMiddleware(), controllers.GetUsersByFolder)
	route.GET("/folder/:id/users", middleware.JWTAuthMiddleware(), controllers.GetSharedUsersForFolder)
	route.GET("/folder/:id/groups", middleware.JWTAuthMiddleware(), controllers.GetSharedGroupsForFolder)
	route.GET("/folder/:id/credential", middleware.JWTAuthMiddleware(), controllers.GetCredentialsByFolder)
	route.GET("/folders/", middleware.JWTAuthMiddleware(), controllers.GetAccessibleFolders)

	route.POST("/share-credentials/users", middleware.JWTAuthMiddleware(), controllers.ShareCredentialsWithUsers)
	route.POST("/share-credentials/groups", middleware.JWTAuthMiddleware(), controllers.ShareCredentialsWithGroups)
	route.POST("share-folder/users", middleware.JWTAuthMiddleware(), controllers.ShareFolderWithUsers)
	route.POST("share-folder/groups", middleware.JWTAuthMiddleware(), controllers.ShareFolderWithGroups)

	// Credential Routes
	route.POST("/credential/", middleware.JWTAuthMiddleware(), controllers.AddCredential)
	route.GET("/credential/:id", middleware.JWTAuthMiddleware(), controllers.FetchCredentialByID)
	route.GET("/credential/:id/sensitive", middleware.JWTAuthMiddleware(), controllers.GetSensitiveFieldsCredentialByID)

	// route.GET("/credential/:id", middleware.JWTAuthMiddleware(), controllers.GetCredentialByID)
	route.POST("/group", middleware.JWTAuthMiddleware(), controllers.CreateGroup)
	route.GET("/group/:groupId", middleware.JWTAuthMiddleware(), controllers.GetGroupMembers)
	route.GET("/groups", middleware.JWTAuthMiddleware(), controllers.GetUserGroups)
	route.POST("/group/members", middleware.JWTAuthMiddleware(), controllers.AddMemberToGroup)
	route.POST("/groups/members", middleware.JWTAuthMiddleware(), controllers.GetUsersOfGroups)
	route.GET("/groups/without-access/:folderId", middleware.JWTAuthMiddleware(), controllers.GetGroupsWithoutAccess)

	route.GET("/credentials/fields/:folderId", middleware.JWTAuthMiddleware(), controllers.GetCredentialsFieldsByFolderID)
	route.POST("/credentials/fields/", middleware.JWTAuthMiddleware(), controllers.GetCredentialsFieldsByIds)
	route.POST("/credentials/by-ids", middleware.JWTAuthMiddleware(), controllers.GetCredentialsByIDs)
	route.GET("/urls", middleware.JWTAuthMiddleware(), controllers.GetAllUrlsForUser)

	route.GET("/group/:groupId/encrypted/", middleware.JWTAuthMiddleware(), controllers.FetchEncryptedValuesByGroupID)

	//Add All route
	//TestRoutes(route)
}
