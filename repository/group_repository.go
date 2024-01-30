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

func FetchCredentialIDsWithGroupAccess(ctx *gin.Context, groupID uuid.UUID) ([]uuid.UUID, error) {
	// doing this because in the table the group_id is nullable
	nullableGroupID := uuid.NullUUID{UUID: groupID, Valid: true}
	credentialIDs, err := database.Store.FetchCredentialIDsWithGroupAccess(ctx, nullableGroupID)
	if err != nil {
		return []uuid.UUID{}, err
	}
	return credentialIDs, nil
}

type AddGroupMemberRepositoryParams struct {
	GroupID           uuid.UUID                        `json:"groupId"`
	MemberID          uuid.UUID                        `json:"memberId"`
	MemberRole        string                           `json:"memberRole"`
	UserEncryptedData []dto.CredentialFieldsForUserDto `json:"encryptedFields"`
}

func AddGroupMember(ctx *gin.Context, payload AddGroupMemberRepositoryParams) error {

	args := db.AddMemberToGroupTransactionParams{
		GroupID:           payload.GroupID,
		UserID:            payload.MemberID,
		MemberRole:        payload.MemberRole,
		UserEncryptedData: payload.UserEncryptedData,
	}

	err := database.Store.AddMemberToGroupTransaction(ctx, args)
	if err != nil {
		return err
	}

	return nil
}

func GetCredentialIDAndTypeWithGroupAccess(ctx *gin.Context, groupID uuid.NullUUID) ([]db.GetCredentialIDAndTypeWithGroupAccessRow, error) {
	return database.Store.GetCredentialIDAndTypeWithGroupAccess(ctx, groupID)

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
