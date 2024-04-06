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
	"manager":      2,
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

	hasAccess, err := repository.HasManageAccessForCredential(ctx, db.HasManageAccessForCredentialParams{
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

	hasAccess, err := repository.HasReadAccessForCredential(ctx, db.HasReadAccessForCredentialParams{
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

func UniqueUsersWithHighestAccessForCredential(userAccess []dto.CredentialUserWithAccess) []dto.CredentialUserWithAccess {

	userAccessMap := map[uuid.UUID]dto.CredentialUserWithAccess{}

	for _, access := range userAccess {
		if _, ok := userAccessMap[access.UserID]; !ok {
			userAccessMap[access.UserID] = access
		} else {
			if CredentialAccessLevels[access.AccessType] > CredentialAccessLevels[userAccessMap[access.UserID].AccessType] {
				userAccessMap[access.UserID] = access

			} else if CredentialAccessLevels[access.AccessType] == CredentialAccessLevels[userAccessMap[access.UserID].AccessType] {
				// incase of multiple access rows with same access level, we choose the one with acquired access
				if access.AccessSource == "acquired" {
					userAccessMap[access.UserID] = access
				}
			}
		}
	}

	uniqueUserAccess := []dto.CredentialUserWithAccess{}
	for _, access := range userAccessMap {
		uniqueUserAccess = append(uniqueUserAccess, access)
	}

	return uniqueUserAccess
}

func GetCredentialUsersWithDirectAccess(ctx *gin.Context, credentialID uuid.UUID, caller uuid.UUID) ([]dto.CredentialUserWithAccess, error) {

	if err := VerifyCredentialReadAccessForUser(ctx, credentialID, caller); err != nil {
		return nil, err
	}

	userAccess, err := repository.GetCredentialUsersWithDirectAccess(ctx, credentialID)
	if err != nil {
		return nil, err
	}

	userAccessObjs := []dto.CredentialUserWithAccess{}
	for _, access := range userAccess {
		userAccessObjs = append(userAccessObjs, dto.CredentialUserWithAccess{
			UserID:       access.UserID,
			Name:         access.Name,
			AccessType:   access.AccessType,
			AccessSource: access.AccessSource,
		})
	}

	uniqueAccessObjs := UniqueUsersWithHighestAccessForCredential(userAccessObjs)

	return uniqueAccessObjs, nil
}

func GetCredentialUsersForDataSync(ctx *gin.Context, credentialID uuid.UUID, caller uuid.UUID) ([]db.GetCredentialUsersForDataSyncRow, error) {

	if err := VerifyCredentialReadAccessForUser(ctx, credentialID, caller); err != nil {
		return nil, err
	}

	return repository.GetCredentialUsersForDataSync(ctx, credentialID)
}

func UniqueGroupsWithHighestAccessForCredential(groupAccess []db.GetCredentialGroupsRow) []db.GetCredentialGroupsRow {

	groupAccessMap := map[uuid.UUID]db.GetCredentialGroupsRow{}

	for _, access := range groupAccess {
		groupID := access.GroupID.UUID
		if _, ok := groupAccessMap[groupID]; !ok {
			groupAccessMap[groupID] = access
		} else {
			if CredentialAccessLevels[access.AccessType] > CredentialAccessLevels[groupAccessMap[groupID].AccessType] {
				groupAccessMap[groupID] = access

			} else if CredentialAccessLevels[access.AccessType] == CredentialAccessLevels[groupAccessMap[groupID].AccessType] {
				// incase of multiple access rows with same access level, we choose the one with acquired access
				if access.AccessSource == "acquired" {
					groupAccessMap[groupID] = access
				}
			}
		}
	}

	uniqueGroupAccess := []db.GetCredentialGroupsRow{}
	for _, access := range groupAccessMap {
		uniqueGroupAccess = append(uniqueGroupAccess, access)
	}

	return uniqueGroupAccess
}

func GetCredentialGroups(ctx *gin.Context, credentialID uuid.UUID, caller uuid.UUID) ([]db.GetCredentialGroupsRow, error) {

	if err := VerifyCredentialReadAccessForUser(ctx, credentialID, caller); err != nil {
		return nil, err
	}

	accessRows, err := repository.GetCredentialGroups(ctx, credentialID)
	if err != nil {
		return nil, err
	}

	uniqueAccessRows := UniqueGroupsWithHighestAccessForCredential(accessRows)

	return uniqueAccessRows, nil
}
