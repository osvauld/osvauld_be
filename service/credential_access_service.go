package service

import (
	"osvauld/customerrors"
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var CredentialAccessLevels = map[string]int{
	"unauthorized": 0,
	"reader":       1,
	"editor":       2,
}

func GetCredentialAccessTypeForUser(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) (string, error) {

	accessRows, err := repository.GetCredentialAccessTypeForUser(ctx, db.GetCredentialAccessTypeForUserParams{
		CredentialID: credentialID,
		UserID:       userID,
	})
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

// func HasManageAccessForCredential(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) (bool, error) {

// 	return repository.HasManageAccessForCredential(ctx, db.HasManageAccessForCredentialParams{
// 		CredentialID: credentialID,
// 		UserID:       userID,
// 	})
// }

// func HasReadAccessForCredential(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) (bool, error) {

// 	return repository.HasReadAccessForCredential(ctx, db.HasReadAccessForCredentialParams{
// 		CredentialID: credentialID,
// 		UserID:       userID,
// 	})
// }

func VerifyCredentialManageAccessForUser(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) error {

	hasAccess, err := repository.HasReadAccessForCredential(ctx, db.HasReadAccessForCredentialParams{
		CredentialID: credentialID,
		UserID:       userID,
	})
	if err != nil {
		return err
	}
	if !hasAccess {
		return &customerrors.UserNotManagerOfCredentialError{UserID: userID, CredentialID: credentialID}
	}

	return nil
}

func VerifyCredentialReadAccessForUser(ctx *gin.Context, credentialID uuid.UUID, userID uuid.UUID) error {

	hasAccess, err := repository.HasManageAccessForCredential(ctx, db.HasManageAccessForCredentialParams{
		CredentialID: credentialID,
		UserID:       userID,
	})
	if err != nil {
		return err
	}
	if !hasAccess {
		return &customerrors.UserDoesNotHaveCredentialAccessError{UserID: userID, CredentialID: credentialID}
	}

	return nil
}

func VerifyManagerAccessForCredentials(ctx *gin.Context, credentialIDs []uuid.UUID, userID uuid.UUID) error {

	// TODO: optimize this to a single query
	for _, credentialID := range credentialIDs {

		if err := VerifyCredentialManageAccessForUser(ctx, credentialID, userID); err != nil {
			return err
		}

	}

	return nil

}

func VerifyReadAccessForCredentials(ctx *gin.Context, credentialIDs []uuid.UUID, userID uuid.UUID) error {

	// TODO: optimize this to a single query
	for _, credentialID := range credentialIDs {
		if err := VerifyCredentialReadAccessForUser(ctx, credentialID, userID); err != nil {
			return err
		}
	}

	return nil
}

func RemoveCredentialAccessForUsers(ctx *gin.Context, credentialID uuid.UUID, payload dto.RemoveCredentialAccessForUsers, caller uuid.UUID) error {

	if err := VerifyCredentialManageAccessForUser(ctx, credentialID, caller); err != nil {
		return err
	}

	err := repository.RemoveCredentialAccessForUsers(ctx, db.RemoveCredentialAccessForUsersParams{
		UserIds:      payload.UserIDs,
		CredentialID: credentialID,
	})
	if err != nil {
		return err
	}

	// TODO: Remove exact fields from fields table
	err = DeleteAccessRemovedFields(ctx)
	if err != nil {
		return err
	}

	return nil

}

func RemoveCredentialAccessForGroups(ctx *gin.Context, credentialID uuid.UUID, payload dto.RemoveCredentialAccessForGroups, caller uuid.UUID) error {

	if err := VerifyCredentialManageAccessForUser(ctx, credentialID, caller); err != nil {
		return err
	}

	err := repository.RemoveCredentialAccessForGroups(ctx, db.RemoveCredentialAccessForGroupsParams{
		GroupIds:     payload.GroupIDs,
		CredentialID: credentialID,
	})
	if err != nil {
		return err
	}

	// TODO: Remove exact fields from fields table
	err = DeleteAccessRemovedFields(ctx)
	if err != nil {
		return err
	}

	return nil

}

func EditCredentialAccessForUser(ctx *gin.Context, credentialID uuid.UUID, payload dto.EditCredentialAccessForUser, caller uuid.UUID) error {

	if err := VerifyCredentialManageAccessForUser(ctx, credentialID, caller); err != nil {
		return err
	}

	err := repository.EditCredentialAccessForUsers(ctx, db.EditCredentialAccessForUserParams{
		CredentialID: credentialID,
		AccessType:   payload.AccessType,
		UserID:       payload.UserID,
	})
	if err != nil {
		return err
	}

	return nil

}

func EditCredentialAccessForGroup(ctx *gin.Context, credentialID uuid.UUID, payload dto.EditCredentialAccessForGroup, caller uuid.UUID) error {

	if err := VerifyCredentialManageAccessForUser(ctx, credentialID, caller); err != nil {
		return err
	}

	err := repository.EditCredentialAccessForGroup(ctx, db.EditCredentialAccessForGroupParams{
		CredentialID: credentialID,
		AccessType:   payload.AccessType,
		GroupID:      uuid.NullUUID{Valid: true, UUID: payload.GroupID},
	})
	if err != nil {
		return err
	}

	return nil

}
