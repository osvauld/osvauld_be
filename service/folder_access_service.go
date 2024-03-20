package service

import (
	"osvauld/customerrors"
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// func HasManageAccessForFolder(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) (bool, error) {

// 	return repository.HasManageAccessForFolder(ctx, db.HasManageAccessForFolderParams{
// 		FolderID: folderID,
// 		UserID:   userID,
// 	})
// }

// func HasReadAccessForFolder(ctx *gin.Context, folderID uuid.UUID, userID uuid.UUID) (bool, error) {

// 	return repository.HasReadAccessForFolder(ctx, db.HasReadAccessForFolderParams{
// 		FolderID: folderID,
// 		UserID:   userID,
// 	})
// }

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
