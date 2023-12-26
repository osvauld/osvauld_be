package service

import (
	"osvauld/customerrors"
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddGroup(ctx *gin.Context, group dto.CreateGroup, userID uuid.UUID) error {
	err := repository.AddGroup(ctx, group, userID)
	return err
}

func AddMembersToGroup(ctx *gin.Context, payload dto.AddMembers, userID uuid.UUID) error {

	err := repository.AddMembersToGroup(ctx, payload, userID)
	return err
}

func GetUserGroups(ctx *gin.Context, userID uuid.UUID) ([]db.Grouping, error) {
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

func FetchEncryptedDataWithGroupAccess(ctx *gin.Context, caller uuid.UUID, groupID uuid.UUID) ([]dto.CredentialEncryptedFielsdDto, error) {

	isMember, err := repository.CheckUserMemberOfGroup(ctx, caller, groupID)
	if !isMember {
		return []dto.CredentialEncryptedFielsdDto{}, &customerrors.UserNotAuthenticatedError{Message: "user does not have access to the group"}
	}
	if err != nil {
		return []dto.CredentialEncryptedFielsdDto{}, err
	}

	credentialIDs, err := repository.FetchCredentialIDsWithGroupAccess(ctx, groupID)
	if err != nil {
		return []dto.CredentialEncryptedFielsdDto{}, err
	}

	allCredentialEncryptedFields := []dto.CredentialEncryptedFielsdDto{}
	for _, credentialID := range credentialIDs {

		credentialEncryptedFields, err := repository.FetchEncryptedFieldsByCredentialIDByAndUserID(ctx, credentialID, caller)
		if err != nil {
			return []dto.CredentialEncryptedFielsdDto{}, err
		}

		dtoObject := dto.CredentialEncryptedFielsdDto{
			CredentialID:    credentialID,
			EncryptedFields: credentialEncryptedFields,
		}

		allCredentialEncryptedFields = append(allCredentialEncryptedFields, dtoObject)
	}

	return allCredentialEncryptedFields, nil
}

func AddMemberToGroup(ctx *gin.Context, payload dto.AddMembers, userID uuid.UUID) error {
	err := repository.AddMembersToGroup(ctx, payload, userID)
	return err
}
