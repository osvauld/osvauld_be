package repository

import (
	dto "osvauld/dtos"
	"osvauld/infra/database"

	"github.com/gin-gonic/gin"
)

func ShareCredentialWithUser(ctx *gin.Context, payload dto.CredentialEncryptedFieldsForUserDto) error {

	err := database.Store.ShareCredentialWithUserTransaction(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}

func ShareCredentialWithGroup(ctx *gin.Context, payload dto.CredentialEncryptedFieldsForGroupDto) error {

	err := database.Store.ShareCredentialWithGroupTransaction(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}

func ShareMultipleCredentialsWithMultipleUsers(ctx *gin.Context, payload []dto.CredentialEncryptedFieldsForUserDto) error {

	err := database.Store.ShareMultipleCredentialsWithMultipleUsersTransaction(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}

func ShareMultipleCredentialsWithMultipleGroups(ctx *gin.Context, payload []dto.CredentialEncryptedFieldsForGroupDto) error {

	err := database.Store.ShareMultipleCredentialsWithMultipleGroupsTransaction(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}
