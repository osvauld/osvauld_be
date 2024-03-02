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

func CheckUserMemberOfGroup(ctx *gin.Context, userID uuid.UUID, groupID uuid.UUID) (bool, error) {
	args := db.CheckUserMemberOfGroupParams{
		UserID:     userID,
		GroupingID: groupID,
	}
	isMember, err := database.Store.CheckUserMemberOfGroup(ctx, args)
	if err != nil {
		return false, err
	}
	return isMember, nil
}

func CheckUserManagerOfGroup(ctx *gin.Context, userID uuid.UUID, groupID uuid.UUID) (bool, error) {
	args := db.FetchGroupAccessTypeParams{
		UserID:     userID,
		GroupingID: groupID,
	}
	role, err := database.Store.FetchGroupAccessType(ctx, args)
	if err != nil {
		return false, err
	}

	if role == "manager" {
		return true, nil
	}

	return false, nil
}

func FetchCredentialIDsWithGroupAccess(ctx *gin.Context, groupID uuid.UUID, caller uuid.UUID) ([]uuid.UUID, error) {
	// doing this because in the table the group_id is nullable
	nullableGroupID := uuid.NullUUID{UUID: groupID, Valid: true}
	credentialIDs, err := database.Store.FetchCredentialIDsWithGroupAccess(ctx, db.FetchCredentialIDsWithGroupAccessParams{GroupID: nullableGroupID, UserID: caller})
	if err != nil {
		return []uuid.UUID{}, err
	}
	return credentialIDs, nil
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
	users, err := database.Store.FetchUsersByGroupIds(ctx, groupIDs)
	if err != nil {
		return []db.FetchUsersByGroupIdsRow{}, err
	}
	return users, nil
}

func GetUsersWithoutGroupAccess(ctx *gin.Context, groupId uuid.UUID) ([]db.GetUsersWithoutGroupAccessRow, error) {
	return database.Store.GetUsersWithoutGroupAccess(ctx, groupId)
}
