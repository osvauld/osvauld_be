package repository

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/database"
	"osvauld/infra/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateGroup(ctx *gin.Context, group dto.CreateGroup, createdBy uuid.UUID) (uuid.UUID, error) {
	arg := db.CreateGroupParams{
		Name:      group.Name,
		CreatedBy: createdBy,
	}
	groupID, err := database.Store.CreateGroupAndAddManager(ctx, arg)
	if err != nil {
		return uuid.UUID{}, err
	}

	return groupID, nil
}

func GetUserGroups(ctx *gin.Context, userID uuid.UUID) ([]dto.GroupDetails, error) {

	userGroups := []dto.GroupDetails{}
	// TODO: This query can return duplicate groups, need to fix this
	groups, err := database.Store.FetchUserGroups(ctx, userID)
	if err != nil {
		logger.Errorf(err.Error())
		return []dto.GroupDetails{}, err
	}

	for _, group := range groups {
		userGroups = append(userGroups, dto.GroupDetails{
			GroupID:   group.ID,
			Name:      group.Name,
			CreatedBy: group.CreatedBy,
			CreatedAt: group.CreatedAt,
		})
	}

	return userGroups, nil
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
	GroupID           uuid.UUID                                 `json:"groupId"`
	MemberID          uuid.UUID                                 `json:"memberId"`
	MemberRole        string                                    `json:"memberRole"`
	UserEncryptedData []dto.CredentialEncryptedFieldsForUserDto `json:"encryptedFields"`
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

func FetchCredentialAccessTypeForGroupMember(ctx *gin.Context, credentialID uuid.UUID, groupID uuid.UUID, userID uuid.UUID) (string, error) {
	args := db.FetchCredentialAccessTypeForGroupMemberParams{
		GroupID:      uuid.NullUUID{UUID: groupID, Valid: true},
		CredentialID: credentialID,
		UserID:       userID,
	}

	accessType, err := database.Store.FetchCredentialAccessTypeForGroupMember(ctx, args)
	if err != nil {
		return "", err
	}

	return accessType, nil
}
