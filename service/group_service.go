package service

import (
	"osvauld/customerrors"
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddGroup(ctx *gin.Context, groupName string, caller uuid.UUID) (dto.GroupDetails, error) {

	groupDto := dto.GroupDetails{
		Name:      groupName,
		CreatedBy: caller,
	}

	groupDetails, err := repository.CreateGroupAndAddManager(ctx, groupDto)

	return groupDetails, err
}

func GetUserGroups(ctx *gin.Context, userID uuid.UUID) ([]db.FetchUserGroupsRow, error) {

	groups, err := repository.GetUserGroups(ctx, userID)
	if err != nil {
		return nil, err
	}

	return groups, nil

}

func GetGroupMembers(ctx *gin.Context, groupID uuid.UUID, userID uuid.UUID) ([]db.GetGroupMembersRow, error) {

	// Check user is a member of the group
	isMember, err := CheckUserMemberOfGroup(ctx, userID, groupID)
	if err != nil {
		return []db.GetGroupMembersRow{}, err
	}
	if !isMember {
		return []db.GetGroupMembersRow{}, &customerrors.UserNotAuthenticatedError{Message: "user is not a member of the group"}
	}

	users, err := repository.GetGroupMembers(ctx, groupID)
	return users, err
}

func CheckUserMemberOfGroup(ctx *gin.Context, userID uuid.UUID, groupID uuid.UUID) (bool, error) {
	isMember, err := repository.CheckUserMemberOfGroup(ctx, userID, groupID)
	return isMember, err
}

func FetchCredentialIDsWithGroupAccess(ctx *gin.Context, caller uuid.UUID, groupID uuid.UUID) ([]uuid.UUID, error) {

	credentialIDs, err := repository.FetchCredentialIDsWithGroupAccess(ctx, groupID)
	return credentialIDs, err
}

func FetchEncryptedDataWithGroupAccess(ctx *gin.Context, caller uuid.UUID, groupID uuid.UUID) ([]dto.CredentialFieldsForUserDto, error) {

	isMember, err := repository.CheckUserMemberOfGroup(ctx, caller, groupID)
	if !isMember {
		return []dto.CredentialFieldsForUserDto{}, &customerrors.UserNotAuthenticatedError{Message: "user does not have access to the group"}
	}
	if err != nil {
		return []dto.CredentialFieldsForUserDto{}, err
	}

	credentialIDs, err := repository.FetchCredentialIDsWithGroupAccess(ctx, groupID)
	if err != nil {
		return []dto.CredentialFieldsForUserDto{}, err
	}

	allCredentialEncryptedFields := []dto.CredentialFieldsForUserDto{}
	for _, credentialID := range credentialIDs {

		credentialEncryptedFields, err := repository.FetchEncryptedFieldsByCredentialIDByAndUserID(ctx, credentialID, caller)
		if err != nil {
			return []dto.CredentialFieldsForUserDto{}, err
		}

		dtoObject := dto.CredentialFieldsForUserDto{
			CredentialID: credentialID,
			Fields:       credentialEncryptedFields,
		}

		allCredentialEncryptedFields = append(allCredentialEncryptedFields, dtoObject)
	}

	return allCredentialEncryptedFields, nil
}

func AddMemberToGroup(ctx *gin.Context, payload dto.AddMemberToGroupRequest, caller uuid.UUID) error {

	isManager, err := repository.CheckUserManagerOfGroup(ctx, caller, payload.GroupID)
	if !isManager {
		return &customerrors.UserNotAuthenticatedError{Message: "caller is not an owner of the group"}
	}
	if err != nil {
		return err
	}

	isMember, err := repository.CheckUserMemberOfGroup(ctx, payload.MemberID, payload.GroupID)
	if isMember {
		return &customerrors.UserAlreadyMemberOfGroupError{Message: "user is already a member of the group"}
	}
	if err != nil {
		return err
	}

	// TODO
	return nil
}

func GetUsersOfGroups(ctx *gin.Context, groupIDs []uuid.UUID) ([]db.FetchUsersByGroupIdsRow, error) {
	users, err := repository.GetUsersOfGroups(ctx, groupIDs)
	return users, err
}
