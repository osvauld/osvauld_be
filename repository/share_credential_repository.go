package repository

import (
	dto "osvauld/dtos"
	"osvauld/infra/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ShareCredentialWithUser(ctx *gin.Context, payload dto.CredentialEncryptedFieldsForUserDto) error {

	err := database.Store.ShareCredentialWithUserTransaction(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}

func ShareCredentialWithGroup(ctx *gin.Context, groupId uuid.UUID, accessType string, encryptedUserData []dto.GroupCredentialPayload) error {

	err := database.Store.ShareCredentialWithGroupTransaction(ctx, groupId, accessType, encryptedUserData)
	if err != nil {
		return err
	}

	return nil
}

func ShareCredentialsWithUsers(ctx *gin.Context, payload []dto.CredentialEncryptedFieldsForUserDto) error {

	err := database.Store.ShareMultipleCredentialsWithMultipleUsersTransaction(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}

func ShareCredentialsWithGroups(ctx *gin.Context, payload []dto.CredentialEncryptedFieldsForGroupDto) error {

	err := database.Store.ShareMultipleCredentialsWithMultipleGroupsTransaction(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}

func ShareFolderWithUsers(ctx *gin.Context, folderId uuid.UUID, payload dto.CredentialsForUsersPayload) error {
	err := database.Store.ShareFolderWithUserTransaction(ctx, folderId, payload)
	if err != nil {
		return err
	}
	return nil
}

func ShareFolderWithGroup(ctx *gin.Context, folderId uuid.UUID, payload dto.CredentialsForGroupsPayload) error {

	err := database.Store.ShareFolderWithGroupTransaction(ctx, folderId, payload)
	if err != nil {
		return err
	}
	return nil
}
