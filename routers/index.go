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
	route.GET("/admin", controllers.GetAdminPage)
	route.POST("/admin", controllers.CreateFirstAdmin)
	route.GET("/health", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"live": "ok"}) })
	route.POST("/user/", controllers.CreateUser)
	route.POST("/user/temp-login", controllers.TempLogin)
	route.DELETE("/user/:id", controllers.RemoveUserFromAll)
	route.POST("/user/register", controllers.Register)
	route.POST("/user/challenge", controllers.GetChallenge)
	route.POST("/user/verify", controllers.VerifyChallenge)
	route.GET("/users", middleware.JWTAuthMiddleware(), controllers.GetAllUsers)
	route.POST("/folder/", middleware.JWTAuthMiddleware(), controllers.CreateFolder)
	route.GET("/folder/:id/credential", middleware.JWTAuthMiddleware(), controllers.GetCredentialsByFolder)
	route.GET("/folders/", middleware.JWTAuthMiddleware(), controllers.FetchAccessibleFoldersForUser)

	route.POST("/share-credentials/users", middleware.JWTAuthMiddleware(), controllers.ShareCredentialsWithUsers)
	route.POST("/share-credentials/groups", middleware.JWTAuthMiddleware(), controllers.ShareCredentialsWithGroups)
	route.POST("share-folder/users", middleware.JWTAuthMiddleware(), controllers.ShareFolderWithUsers)
	route.POST("share-folder/groups", middleware.JWTAuthMiddleware(), controllers.ShareFolderWithGroups)

	route.POST("/credential/", middleware.JWTAuthMiddleware(), controllers.AddCredential)
	route.GET("/credential/:id", middleware.JWTAuthMiddleware(), controllers.GetCredentialDataByID)
	route.PUT("/credential/:id", middleware.JWTAuthMiddleware(), controllers.EditCredential)
	route.GET("/credential/:id/sensitive", middleware.JWTAuthMiddleware(), controllers.GetSensitiveFieldsByCredentialID)

	route.GET("/credential/:id/users-data-sync", middleware.JWTAuthMiddleware(), controllers.GetCredentialUsersForDataSync)
	route.GET("/credential/:id/groups", middleware.JWTAuthMiddleware(), controllers.GetCredentialGroups)
	route.GET("/credential/:id/users", middleware.JWTAuthMiddleware(), controllers.GetCredentialUsersWithDirectAccess)

	route.GET("/folder/:id/users-data-sync", middleware.JWTAuthMiddleware(), controllers.GetFolderUsersForDataSync)
	route.GET("/folder/:id/users", middleware.JWTAuthMiddleware(), controllers.GetFolderUsersWithDirectAccess)
	route.GET("/folder/:id/groups", middleware.JWTAuthMiddleware(), controllers.GetFolderGroups)
	route.DELETE("/folder/:id", middleware.JWTAuthMiddleware(), controllers.RemoveFolder)

	route.POST("/group", middleware.JWTAuthMiddleware(), controllers.CreateGroup)
	route.GET("/group/:groupId", middleware.JWTAuthMiddleware(), controllers.GetGroupMembers)

	route.GET("/groups", middleware.JWTAuthMiddleware(), controllers.GetUserGroups)
	route.POST("/group/members", middleware.JWTAuthMiddleware(), controllers.AddMemberToGroup)
	// TODO: maybe change the route to include groupId and memberId
	route.DELETE("/group/member", middleware.JWTAuthMiddleware(), controllers.RemoveMemberFromGroup)
	route.DELETE("/group/:groupId", middleware.JWTAuthMiddleware(), controllers.RemoveGroup)
	route.POST("/groups/members", middleware.JWTAuthMiddleware(), controllers.GetUsersOfGroups)
	route.GET("/groups/without-access/:folderId", middleware.JWTAuthMiddleware(), controllers.GetGroupsWithoutAccess)

	route.GET("/credentials/fields/:folderId", middleware.JWTAuthMiddleware(), controllers.GetCredentialsFieldsByFolderID)
	route.POST("/credentials/fields/", middleware.JWTAuthMiddleware(), controllers.GetCredentialsFieldsByIds)
	route.POST("/credentials/by-ids", middleware.JWTAuthMiddleware(), controllers.GetCredentialsByIDs)
	route.GET("/credentials/search", middleware.JWTAuthMiddleware(), controllers.GetSearchData)
	route.GET("/urls", middleware.JWTAuthMiddleware(), controllers.GetAllUrlsForUser)

	route.GET("/group/:groupId/credential-fields", middleware.JWTAuthMiddleware(), controllers.GetAllCredentialFieldsByGroupID)
	route.GET("/groups/:groupId/users/without-access", middleware.JWTAuthMiddleware(), controllers.GetUsersWithoutGroupAccess)

	route.POST("/credential/:id/remove-user-access", middleware.JWTAuthMiddleware(), controllers.RemoveCredentialAccessForUsers)
	route.POST("/credential/:id/remove-group-access", middleware.JWTAuthMiddleware(), controllers.RemoveCredentialAccessForGroups)
	route.DELETE("/credential/:id", middleware.JWTAuthMiddleware(), controllers.RemoveCredential)

	route.POST("/folder/:id/remove-user-access", middleware.JWTAuthMiddleware(), controllers.RemoveFolderAccessForUsers)
	route.POST("/folder/:id/remove-group-access", middleware.JWTAuthMiddleware(), controllers.RemoveFolderAccessForGroups)

	route.POST("/credential/:id/edit-user-access", middleware.JWTAuthMiddleware(), controllers.EditCredentialAccessForUser)
	route.POST("/credential/:id/edit-group-access", middleware.JWTAuthMiddleware(), controllers.EditCredentialAccessForGroup)

	route.POST("/folder/:id/edit-user-access", middleware.JWTAuthMiddleware(), controllers.EditFolderAccessForUser)
	route.POST("/folder/:id/edit-group-access", middleware.JWTAuthMiddleware(), controllers.EditFolderAccessForGroup)

}
