package service

import (
	"osvauld/customerrors"
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddGroup(ctx *gin.Context, group dto.CreateGroup, userID uuid.UUID) (uuid.UUID, error) {
	groupID, err := repository.CreateGroup(ctx, group, userID)
	return groupID, err
}

func GetUserGroups(ctx *gin.Context, userID uuid.UUID) ([]dto.GroupDetails, error) {
	groups, err := repository.GetUserGroups(ctx, userID)
	return groups, err

}

func GetGroupMembers(ctx *gin.Context, userID uuid.UUID, groupId uuid.UUID) ([]db.GetGroupMembersRow, error) {
	users, err := repository.GetGroupMembers(ctx, groupId)
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

func FetchEncryptedDataWithGroupAccess(ctx *gin.Context, caller uuid.UUID, groupID uuid.UUID) ([]dto.CredentialEncryptedFieldsForUserDto, error) {

	isMember, err := repository.CheckUserMemberOfGroup(ctx, caller, groupID)
	if !isMember {
		return []dto.CredentialEncryptedFieldsForUserDto{}, &customerrors.UserNotAuthenticatedError{Message: "user does not have access to the group"}
	}
	if err != nil {
		return []dto.CredentialEncryptedFieldsForUserDto{}, err
	}

	credentialIDs, err := repository.FetchCredentialIDsWithGroupAccess(ctx, groupID)
	if err != nil {
		return []dto.CredentialEncryptedFieldsForUserDto{}, err
	}

	allCredentialEncryptedFields := []dto.CredentialEncryptedFieldsForUserDto{}
	for _, credentialID := range credentialIDs {

		credentialEncryptedFields, err := repository.FetchEncryptedFieldsByCredentialIDByAndUserID(ctx, credentialID, caller)
		if err != nil {
			return []dto.CredentialEncryptedFieldsForUserDto{}, err
		}

		dtoObject := dto.CredentialEncryptedFieldsForUserDto{
			CredentialID:    credentialID,
			EncryptedFields: credentialEncryptedFields,
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

	userEncryptedDataWithAccessType := []dto.CredentialEncryptedFieldsForUserDto{}

	for _, credential := range payload.EncryptedData {

		credentialAccessTypeForGroup, err := repository.FetchCredentialAccessTypeForGroupMember(ctx, credential.CredentialID, payload.GroupID, caller)
		if err != nil {
			return err
		}
		// find out the current group access of each credential
		encryptedDataWithAccessType := dto.CredentialEncryptedFieldsForUserDto{
			CredentialID:    credential.CredentialID,
			AccessType:      credentialAccessTypeForGroup,
			EncryptedFields: credential.EncryptedFields,
		}

		userEncryptedDataWithAccessType = append(userEncryptedDataWithAccessType, encryptedDataWithAccessType)

	}

	args := repository.AddGroupMemberRepositoryParams{
		GroupID:           payload.GroupID,
		MemberID:          payload.MemberID,
		MemberRole:        payload.MemberRole,
		UserEncryptedData: userEncryptedDataWithAccessType,
	}

	err = repository.AddGroupMember(ctx, args)
	if err != nil {
		return err
	}

	return nil
}
