package service

import (
	"osvauld/customerrors"
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/logger"
	"osvauld/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func VerifyMemberOfGroup(ctx *gin.Context, groupID uuid.UUID, caller uuid.UUID) error {
	isMember, err := repository.CheckUserMemberOfGroup(ctx, db.CheckUserMemberOfGroupParams{
		UserID:     caller,
		GroupingID: groupID,
	})
	if err != nil {
		return err
	}
	if !isMember {
		return &customerrors.UserNotMemberOfGroupError{UserID: caller, GroupID: groupID}
	}

	return nil
}

func VerifyAdminOfGroup(ctx *gin.Context, groupID uuid.UUID, caller uuid.UUID) error {
	isAdmin, err := repository.CheckUserAdminOfGroup(ctx, db.CheckUserAdminOfGroupParams{
		UserID:     caller,
		GroupingID: groupID,
	})
	if err != nil {
		return err
	}
	if !isAdmin {
		return &customerrors.UserNotAdminOfGroupError{UserID: caller, GroupID: groupID}
	}

	return nil
}

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

func GetGroupMembers(ctx *gin.Context, groupID uuid.UUID, caller uuid.UUID) ([]db.GetGroupMembersRow, error) {

	if err := VerifyMemberOfGroup(ctx, groupID, caller); err != nil {
		return nil, err
	}

	users, err := repository.GetGroupMembers(ctx, groupID)
	return users, err
}

func FetchCredentialIDsWithGroupAccess(ctx *gin.Context, caller uuid.UUID, groupID uuid.UUID) ([]uuid.UUID, error) {

	if err := VerifyMemberOfGroup(ctx, groupID, caller); err != nil {
		return nil, err
	}

	credentialIDs, err := repository.FetchCredentialIDsWithGroupAccess(ctx, db.FetchCredentialIDsWithGroupAccessParams{
		GroupID: uuid.NullUUID{UUID: groupID, Valid: true},
		UserID:  caller,
	})
	return credentialIDs, err
}

func GetCredentialFieldsByGroupID(ctx *gin.Context, caller uuid.UUID, groupID uuid.UUID) ([]dto.CredentialFields, error) {

	if err := VerifyMemberOfGroup(ctx, groupID, caller); err != nil {
		return nil, err
	}

	credentialIDs, err := repository.FetchCredentialIDsWithGroupAccess(ctx, db.FetchCredentialIDsWithGroupAccessParams{
		GroupID: uuid.NullUUID{UUID: groupID, Valid: true},
		UserID:  caller,
	})
	if err != nil {
		return nil, err
	}

	credentialFields, err := GetFieldsByCredentialIDs(ctx, credentialIDs, caller)
	if err != nil {
		return nil, err
	}

	return credentialFields, nil
}

func AddMemberToGroup(ctx *gin.Context, payload dto.AddMemberToGroupRequest, caller uuid.UUID) error {

	if err := VerifyAdminOfGroup(ctx, payload.GroupID, caller); err != nil {
		return err
	}

	credentialIDAndTypeWithGroupAccess, err := repository.GetCredentialAccessDetailsWithGroupAccess(ctx, uuid.NullUUID{UUID: payload.GroupID, Valid: true})
	if err != nil {
		return err
	}

	folderIDAndTypeWithGroupAccess, err := repository.GetFolderIDAndTypeWithGroupAccess(ctx, uuid.NullUUID{UUID: payload.GroupID, Valid: true})
	if err != nil {
		return err
	}

	userFieldRecords := []db.AddFieldParams{}
	for _, credential := range payload.Credentials {

		fieldDataExists, err := repository.CheckFieldEntryExists(ctx, db.CheckFieldEntryExistsParams{
			UserID:       payload.MemberID,
			CredentialID: credential.CredentialID,
		})
		if err != nil {
			return err
		}

		if !fieldDataExists {
			userFields, err := CreateFieldDataRecords(ctx, credential, payload.MemberID, caller)
			if err != nil {
				return err
			}

			userFieldRecords = append(userFieldRecords, userFields...)

		} else {
			logger.Infof("Field data already exists for credential %s and user %s", credential.CredentialID, payload.MemberID)
		}

	}

	credentialAccessRecords := []db.AddCredentialAccessParams{}
	for _, credentialDetails := range credentialIDAndTypeWithGroupAccess {
		credentialAccessRecord := db.AddCredentialAccessParams{
			CredentialID: credentialDetails.CredentialID,
			UserID:       payload.MemberID,
			AccessType:   credentialDetails.AccessType,
			GroupID:      uuid.NullUUID{UUID: payload.GroupID, Valid: true},
			FolderID:     credentialDetails.FolderID,
		}
		credentialAccessRecords = append(credentialAccessRecords, credentialAccessRecord)
	}

	folderAccessRecords := []db.AddFolderAccessParams{}
	for _, folderDetails := range folderIDAndTypeWithGroupAccess {
		folderAccessRecord := db.AddFolderAccessParams{
			FolderID:   folderDetails.FolderID,
			UserID:     payload.MemberID,
			AccessType: folderDetails.AccessType,
			GroupID:    uuid.NullUUID{UUID: payload.GroupID, Valid: true},
		}
		folderAccessRecords = append(folderAccessRecords, folderAccessRecord)
	}

	groupMembershipRecords := []db.AddGroupMemberParams{
		{
			GroupingID: payload.GroupID,
			UserID:     payload.MemberID,
			AccessType: payload.MemberRole,
		},
	}

	addMemberToGroupTransactionParams := db.AddMembersToGroupTransactionParams{
		FieldArgs:            userFieldRecords,
		CredentialAccessArgs: credentialAccessRecords,
		FolderAccessArgs:     folderAccessRecords,
		GroupMembershipArgs:  groupMembershipRecords,
	}

	err = repository.AddMembersToGroupTransaction(ctx, addMemberToGroupTransactionParams)
	if err != nil {
		return err
	}

	return nil
}

func GetUsersWithoutGroupAccess(ctx *gin.Context, groupID uuid.UUID, caller uuid.UUID) ([]db.GetUsersWithoutGroupAccessRow, error) {

	if err := VerifyMemberOfGroup(ctx, groupID, caller); err != nil {
		return nil, err
	}

	users, err := repository.GetUsersWithoutGroupAccess(ctx, groupID)
	if err != nil {
		return []db.GetUsersWithoutGroupAccessRow{}, err
	}
	return users, err
}

func GetUsersOfGroups(ctx *gin.Context, groupIDs []uuid.UUID, caller uuid.UUID) ([]db.FetchUsersByGroupIdsRow, error) {

	for _, groupID := range groupIDs {
		if err := VerifyMemberOfGroup(ctx, groupID, caller); err != nil {
			return nil, err
		}
	}

	users, err := repository.GetUsersOfGroups(ctx, groupIDs)
	return users, err
}

func GetGroupsWithoutAccess(ctx *gin.Context, folderID uuid.UUID, caller uuid.UUID) ([]db.GetGroupsWithoutAccessRow, error) {

	if err := VerifyFolderReadAccessForUser(ctx, folderID, caller); err != nil {
		return nil, err
	}

	groups, err := repository.GetGroupsWithoutAccess(ctx, folderID, caller)
	return groups, err
}

func RemoveMemberFromGroup(ctx *gin.Context, payload dto.RemoveMemberFromGroupRequest, caller uuid.UUID) error {
	//TODO: check the caller privilage

	return repository.RemoveUserFromGroupList(ctx, payload.MemberID, payload.GroupID)
}

func RemoveGroup(ctx *gin.Context, groupID uuid.UUID, caller uuid.UUID) error {
	// TODO: check the caller privilage
	return repository.RemoveGroup(ctx, groupID)
}
