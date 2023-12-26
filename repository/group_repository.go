package repository

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/database"
	"osvauld/infra/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddGroup(ctx *gin.Context, group dto.CreateGroup, userID uuid.UUID) error {
	arg := db.CreateGroupParams{
		Name:   group.Name,
		UserID: userID,
	}
	err := database.Store.CreateGroup(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	return nil
}

func AddMembersToGroup(ctx *gin.Context, payload dto.AddMembers, userID uuid.UUID) error {
	// uuidArray := pq.Array(payload.Members)
	arg := db.AddMemberToGroupParams{
		GroupingID: payload.GroupID,
		Column2:    payload.Members,
	}
	err := database.Store.AddMemberToGroup(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	return nil

}

func GetUserGroups(ctx *gin.Context, userID uuid.UUID) ([]db.Grouping, error) {

	groups, err := database.Store.GetUserGroups(ctx, userID)
	if err != nil {
		logger.Errorf(err.Error())
		return groups, err
	}
	return groups, nil
}

func GetGroupMembers(ctx *gin.Context, groupID uuid.UUID) ([]db.GetGroupMembersRow, error) {
	users, err := database.Store.GetGroupMembers(ctx, groupID)
	if err != nil {
		logger.Errorf(err.Error())
		return users, err
	}
	return users, nil
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

func FetchCredentialIDsWithGroupAccess(ctx *gin.Context, groupID uuid.UUID) ([]uuid.UUID, error) {
	// doing this because in the table the group_id is nullable
	nullableGroupID := uuid.NullUUID{UUID: groupID, Valid: true}
	credentialIDs, err := database.Store.FetchCredentialIDsWithGroupAccess(ctx, nullableGroupID)
	if err != nil {
		return []uuid.UUID{}, err
	}
	return credentialIDs, nil
}
