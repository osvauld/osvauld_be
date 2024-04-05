package repository

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateGroupAndAddManager(ctx *gin.Context, groupData dto.GroupDetails) (dto.GroupDetails, error) {
	return database.Store.CreateGroupAndAddManager(ctx, groupData)
}

func GetUserGroups(ctx *gin.Context, userID uuid.UUID) ([]db.FetchUserGroupsRow, error) {
	return database.Store.FetchUserGroups(ctx, userID)
}

func GetGroupMembers(ctx *gin.Context, groupID uuid.UUID) ([]db.GetGroupMembersRow, error) {
	return database.Store.GetGroupMembers(ctx, groupID)
}

func CheckUserMemberOfGroup(ctx *gin.Context, args db.CheckUserMemberOfGroupParams) (bool, error) {
	return database.Store.CheckUserMemberOfGroup(ctx, args)
}

func CheckUserAdminOfGroup(ctx *gin.Context, args db.CheckUserAdminOfGroupParams) (bool, error) {
	return database.Store.CheckUserAdminOfGroup(ctx, args)
}

func FetchCredentialIDsWithGroupAccess(ctx *gin.Context, args db.FetchCredentialIDsWithGroupAccessParams) ([]uuid.UUID, error) {
	return database.Store.FetchCredentialIDsWithGroupAccess(ctx, args)
}

func AddMembersToGroupTransaction(ctx *gin.Context, args db.AddMembersToGroupTransactionParams) error {
	return database.Store.AddMembersToGroupTransaction(ctx, args)
}

func GetCredentialAccessDetailsWithGroupAccess(ctx *gin.Context, groupID uuid.NullUUID) ([]db.GetCredentialAccessDetailsWithGroupAccessRow, error) {
	return database.Store.GetCredentialAccessDetailsWithGroupAccess(ctx, groupID)
}

func GetFolderIDAndTypeWithGroupAccess(ctx *gin.Context, groupID uuid.NullUUID) ([]db.GetFolderIDAndTypeWithGroupAccessRow, error) {
	return database.Store.GetFolderIDAndTypeWithGroupAccess(ctx, groupID)
}

func GetUsersOfGroups(ctx *gin.Context, groupIDs []uuid.UUID) ([]db.FetchUsersByGroupIdsRow, error) {
	return database.Store.FetchUsersByGroupIds(ctx, groupIDs)
}

func GetUsersWithoutGroupAccess(ctx *gin.Context, groupId uuid.UUID) ([]db.GetUsersWithoutGroupAccessRow, error) {
	return database.Store.GetUsersWithoutGroupAccess(ctx, groupId)
}

func RemoveMemberFromGroup(ctx *gin.Context, args db.RemoveMemberFromGroupTransactionParams) error {
	return database.Store.RemoveMemberFromGroupTransaction(ctx, args)
}

func RemoveGroup(ctx *gin.Context, groupID uuid.UUID) error {
	return database.Store.RemoveGroup(ctx, groupID)
}
