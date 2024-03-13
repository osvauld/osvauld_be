package service

import (
	"fmt"
	"osvauld/customerrors"
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/logger"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var CredentialAccessLevels = map[string]int{
	"unauthorized": 0,
	"read":         1,
	"write":        2,
	"owner":        99,
}

func GetAccessTypeForCredential(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) (string, error) {

	accessRows, err := repository.GetCredentialAccessForUser(ctx, credentialID, userID)
	if err != nil {
		return "", err
	}

	// Find the higest access level for a user
	highestAccess := "unauthorized"
	for _, accessRow := range accessRows {

		if CredentialAccessLevels[accessRow.AccessType] > CredentialAccessLevels[highestAccess] {
			highestAccess = accessRow.AccessType
		}
	}

	return highestAccess, nil

}

func HasReadAccessForCredential(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) (bool, error) {
	access, err := GetAccessTypeForCredential(ctx, credentialID, userID)
	logger.Debugf("access level for credential %s, user: %s is %s", credentialID, userID, access)
	if err != nil {
		return false, err
	}
	if CredentialAccessLevels[access] > 0 {
		return true, nil
	}
	return false, nil
}

func CheckHasReadAccessForCredential(ctx *gin.Context, access string) bool {
	return CredentialAccessLevels[access] > 0

}

func HasOwnerAccessForFolder(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) (bool, error) {

	return repository.HasOwnerAccessForFolder(ctx, &db.HasOwnerAccessForFolderParams{
		UserID:   userID,
		FolderID: folderID,
	})
}

func HasOwnerAccessForCredential(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) (bool, error) {

	access, err := GetAccessTypeForCredential(ctx, credentialID, userID)
	logger.Debugf("access level for credential %s, user: %s is %s", credentialID, userID, access)
	if err != nil {
		return false, err
	}
	if CredentialAccessLevels[access] == 99 {
		return true, nil
	}
	return false, nil

}

func HasWriteAccessForCredential(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) (bool, error) {
	access, err := GetAccessTypeForCredential(ctx, credentialID, userID)
	if err != nil {
		return false, err
	}
	if CredentialAccessLevels[access] > 1 {
		return true, nil
	}
	return false, nil
}

func HasOwnerAccessForCredentials(ctx *gin.Context, credentialIDs []uuid.UUID, userID uuid.UUID) (bool, error) {

	// TODO: optimize this
	for _, credentialID := range credentialIDs {

		isOwner, err := HasOwnerAccessForCredential(ctx, credentialID, userID)
		if err != nil {
			return false, err
		}
		if !isOwner {
			errMsg := fmt.Sprintf("user %s does not have owner access for credential %s", userID, credentialID)
			return false, &customerrors.UserNotAnOwnerOfCredentialError{Message: errMsg}
		}
	}

	return true, nil

}

func RemoveCredentialAccessForUsers(ctx *gin.Context, credentialID uuid.UUID, payload dto.RemoveCredentialAccessForUsers, caller uuid.UUID) error {

	// Check caller has owner access for credential
	isOwner, err := HasOwnerAccessForCredentials(ctx, []uuid.UUID{credentialID}, caller)
	if err != nil {
		return err
	}

	if !isOwner {
		errMsg := fmt.Sprintf("user %s does not have owner access for credential %s", caller, credentialID)
		return &customerrors.UserNotAnOwnerOfCredentialError{Message: errMsg}
	}

	err = repository.RemoveCredentialAccessForUsers(ctx, db.RemoveCredentialAccessForUsersParams{
		UserIds:      payload.UserIDs,
		CredentialID: credentialID,
	})
	if err != nil {
		return err
	}

	// TODO: Remove extact fields from fields table
	err = DeleteAccessRemovedFields(ctx)
	if err != nil {
		return err
	}

	return nil

}

func RemoveFolderAccessForUsers(ctx *gin.Context, folderID uuid.UUID, payload dto.RemoveFolderAccessForUsers, caller uuid.UUID) error {

	// Check caller has owner access for folder
	isOwner, err := HasOwnerAccessForFolder(ctx, folderID, caller)
	if err != nil {
		return err
	}

	if !isOwner {
		errMsg := fmt.Sprintf("user %s does not have owner access for folder %s", caller, folderID)
		return &customerrors.UserNotAnOwnerOfFolderError{Message: errMsg}
	}

	err = repository.RemoveFolderAccessForUser(ctx, db.RemoveFolderAccessForUsersParams{
		UserIds:  payload.UserIDs,
		FolderID: folderID,
	})
	if err != nil {
		return err
	}

	// TODO: Remove extact fields from fields table
	err = DeleteAccessRemovedFields(ctx)
	if err != nil {
		return err
	}

	return nil

}

func RemoveCredentialAccessForGroups(ctx *gin.Context, credentialID uuid.UUID, payload dto.RemoveCredentialAccessForGroups, caller uuid.UUID) error {

	// Check caller has owner access for credential
	isOwner, err := HasOwnerAccessForCredentials(ctx, []uuid.UUID{credentialID}, caller)
	if err != nil {
		return err
	}

	if !isOwner {
		errMsg := fmt.Sprintf("user %s does not have owner access for credential %s", caller, credentialID)
		return &customerrors.UserNotAnOwnerOfCredentialError{Message: errMsg}
	}

	err = repository.RemoveCredentialAccessForGroups(ctx, db.RemoveCredentialAccessForGroupsParams{
		GroupIds:     payload.GroupIDs,
		CredentialID: credentialID,
	})
	if err != nil {
		return err
	}

	// TODO: Remove extact fields from fields table
	err = DeleteAccessRemovedFields(ctx)
	if err != nil {
		return err
	}

	return nil

}

func RemoveFolderAccessForGroups(ctx *gin.Context, folderID uuid.UUID, payload dto.RemoveFolderAccessForGroups, caller uuid.UUID) error {

	// Check caller has owner access for folder
	isOwner, err := HasOwnerAccessForFolder(ctx, folderID, caller)
	if err != nil {
		return err
	}

	if !isOwner {
		errMsg := fmt.Sprintf("user %s does not have owner access for folder %s", caller, folderID)
		return &customerrors.UserNotAnOwnerOfFolderError{Message: errMsg}
	}

	err = repository.RemoveFolderAccessForGroups(ctx, db.RemoveFolderAccessForGroupsParams{
		GroupIds: payload.GroupIDs,
		FolderID: folderID,
	})
	if err != nil {
		return err
	}

	// TODO: Remove extact fields from fields table
	err = DeleteAccessRemovedFields(ctx)
	if err != nil {
		return err
	}

	return nil

}

func EditCredentialAccessForUser(ctx *gin.Context, credentialID uuid.UUID, payload dto.EditCredentialAccessForUser, caller uuid.UUID) error {

	// Check caller has owner access for credential
	isOwner, err := HasOwnerAccessForCredentials(ctx, []uuid.UUID{credentialID}, caller)
	if err != nil {
		return err
	}

	if !isOwner {
		errMsg := fmt.Sprintf("user %s does not have owner access for credential %s", caller, credentialID)
		return &customerrors.UserNotAnOwnerOfCredentialError{Message: errMsg}
	}

	err = repository.EditCredentialAccessForUsers(ctx, db.EditCredentialAccessForUserParams{
		CredentialID: credentialID,
		AccessType:   payload.AccessType,
		UserID:       payload.UserID,
	})
	if err != nil {
		return err
	}

	return nil

}

func EditFolderAccessForUser(ctx *gin.Context, folderID uuid.UUID, payload dto.EditFolderAccessForUser, caller uuid.UUID) error {

	// Check caller has owner access for folder
	isOwner, err := HasOwnerAccessForFolder(ctx, folderID, caller)
	if err != nil {
		return err
	}

	if !isOwner {
		errMsg := fmt.Sprintf("user %s does not have owner access for folder %s", caller, folderID)
		return &customerrors.UserNotAnOwnerOfFolderError{Message: errMsg}
	}

	err = repository.EditFolderAccessForUser(ctx, db.EditFolderAccessForUserParams{
		FolderID:   folderID,
		AccessType: payload.AccessType,
		UserID:     payload.UserID,
	})
	if err != nil {
		return err
	}

	return nil

}

func EditCredentialAccessForGroup(ctx *gin.Context, credentialID uuid.UUID, payload dto.EditCredentialAccessForGroup, caller uuid.UUID) error {

	// Check caller has owner access for credential
	isOwner, err := HasOwnerAccessForCredentials(ctx, []uuid.UUID{credentialID}, caller)
	if err != nil {
		return err
	}

	if !isOwner {
		errMsg := fmt.Sprintf("user %s does not have owner access for credential %s", caller, credentialID)
		return &customerrors.UserNotAnOwnerOfCredentialError{Message: errMsg}
	}

	err = repository.EditCredentialAccessForGroup(ctx, db.EditCredentialAccessForGroupParams{
		CredentialID: credentialID,
		AccessType:   payload.AccessType,
		GroupID:      uuid.NullUUID{Valid: true, UUID: payload.GroupID},
	})
	if err != nil {
		return err
	}

	return nil

}

func EditFolderAccessForGroup(ctx *gin.Context, folderID uuid.UUID, payload dto.EditFolderAccessForGroup, caller uuid.UUID) error {

	// Check caller has owner access for folder
	isOwner, err := HasOwnerAccessForFolder(ctx, folderID, caller)
	if err != nil {
		return err
	}

	if !isOwner {
		errMsg := fmt.Sprintf("user %s does not have owner access for folder %s", caller, folderID)
		return &customerrors.UserNotAnOwnerOfFolderError{Message: errMsg}
	}

	err = repository.EditFolderAccessForGroup(ctx, db.EditFolderAccessForGroupParams{
		FolderID:   folderID,
		AccessType: payload.AccessType,
		GroupID:    uuid.NullUUID{Valid: true, UUID: payload.GroupID},
	})
	if err != nil {
		return err
	}

	return nil

}
