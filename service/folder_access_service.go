package service

import (
	"osvauld/customerrors"
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var folderAccessLevels = map[string]int{
	"unauthorized": 0,
	"reader":       1,
	"manager":      99,
}

func VerifyFolderManageAccessForUser(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) error {

	hasAccess, err := repository.HasManageAccessForFolder(ctx, db.HasManageAccessForFolderParams{
		FolderID: folderID,
		UserID:   userID,
	})
	if err != nil {
		return err
	}

	if !hasAccess {
		return &customerrors.UserNotManagerOfFolderError{UserID: userID, FolderID: folderID}
	}

	return nil
}

func VerifyFolderReadAccessForUser(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) error {

	hasAccess, err := repository.HasReadAccessForFolder(ctx, db.HasReadAccessForFolderParams{
		FolderID: folderID,
		UserID:   userID,
	})
	if err != nil {
		return err
	}

	if !hasAccess {
		return &customerrors.UserDoesNotHaveFolderAccessError{UserID: userID, FolderID: folderID}
	}

	return nil
}

func RemoveFolderAccessForUsers(ctx *gin.Context, folderID uuid.UUID, payload dto.RemoveFolderAccessForUsers, caller uuid.UUID) error {

	if err := VerifyFolderManageAccessForUser(ctx, folderID, caller); err != nil {
		return err
	}

	err := repository.RemoveFolderAccessForUser(ctx, db.RemoveFolderAccessForUsersParams{
		UserIds:  payload.UserIDs,
		FolderID: folderID,
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

func RemoveFolderAccessForGroups(ctx *gin.Context, folderID uuid.UUID, payload dto.RemoveFolderAccessForGroups, caller uuid.UUID) error {

	if err := VerifyFolderManageAccessForUser(ctx, folderID, caller); err != nil {
		return err
	}

	err := repository.RemoveFolderAccessForGroups(ctx, db.RemoveFolderAccessForGroupsParams{
		GroupIds: payload.GroupIDs,
		FolderID: folderID,
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

func EditFolderAccessForUser(ctx *gin.Context, folderID uuid.UUID, payload dto.EditFolderAccessForUser, caller uuid.UUID) error {

	if err := VerifyFolderManageAccessForUser(ctx, folderID, caller); err != nil {
		return err
	}

	err := repository.EditFolderAccessForUser(ctx, db.EditFolderAccessForUserParams{
		FolderID:   folderID,
		AccessType: payload.AccessType,
		UserID:     payload.UserID,
	})
	if err != nil {
		return err
	}

	return nil

}

func EditFolderAccessForGroup(ctx *gin.Context, folderID uuid.UUID, payload dto.EditFolderAccessForGroup, caller uuid.UUID) error {

	if err := VerifyFolderManageAccessForUser(ctx, folderID, caller); err != nil {
		return err
	}

	err := repository.EditFolderAccessForGroup(ctx, db.EditFolderAccessForGroupParams{
		FolderID:   folderID,
		AccessType: payload.AccessType,
		GroupID:    uuid.NullUUID{Valid: true, UUID: payload.GroupID},
	})
	if err != nil {
		return err
	}

	return nil

}

func UniqueUsersWithHighestAccessForFolder(userAccess []dto.FolderUserWithAccess) []dto.FolderUserWithAccess {

	userAccessMap := map[uuid.UUID]dto.FolderUserWithAccess{}

	for _, access := range userAccess {
		if _, ok := userAccessMap[access.UserID]; !ok {
			userAccessMap[access.UserID] = access
		} else if CredentialAccessLevels[access.AccessType] > CredentialAccessLevels[userAccessMap[access.UserID].AccessType] {
			userAccessMap[access.UserID] = access

		}
	}

	uniqueUserAccess := []dto.FolderUserWithAccess{}
	for _, access := range userAccessMap {
		uniqueUserAccess = append(uniqueUserAccess, access)
	}

	return uniqueUserAccess

}

func GetFolderUsersWithDirectAccess(ctx *gin.Context, folderID uuid.UUID, caller uuid.UUID) ([]dto.FolderUserWithAccess, error) {

	if err := VerifyFolderReadAccessForUser(ctx, folderID, caller); err != nil {
		return nil, err
	}

	userAccess, err := repository.GetFolderUsersWithDirectAccess(ctx, folderID)
	if err != nil {
		return nil, err
	}

	userAccessObjs := []dto.FolderUserWithAccess{}
	for _, access := range userAccess {
		userAccessObjs = append(userAccessObjs, dto.FolderUserWithAccess{
			UserID:     access.UserID,
			Name:       access.Name,
			AccessType: access.AccessType,
		})
	}

	userAccessObjs = UniqueUsersWithHighestAccessForFolder(userAccessObjs)

	return userAccessObjs, nil

}

func GetFolderUsersForDataSync(ctx *gin.Context, folderID uuid.UUID, caller uuid.UUID) ([]db.GetFolderUsersForDataSyncRow, error) {

	if err := VerifyFolderReadAccessForUser(ctx, folderID, caller); err != nil {
		return nil, err
	}

	return repository.GetFolderUsersForDataSync(ctx, folderID)

}

func GetFolderGroups(ctx *gin.Context, folderID uuid.UUID, caller uuid.UUID) ([]db.GetFolderGroupsRow, error) {

	if err := VerifyFolderReadAccessForUser(ctx, folderID, caller); err != nil {
		return nil, err
	}

	return repository.GetFolderGroups(ctx, folderID)

}
