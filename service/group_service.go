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

func GetGroupMembers(ctx *gin.Context, groupID uuid.UUID, userID uuid.UUID) ([]db.GetGroupMembersRow, error) {

	// Check user is a member of the group
	isMember, err := CheckUserMemberOfGroup(ctx, userID, groupID)
	if err != nil {
		return []db.GetGroupMembersRow{}, err
	}
	if !isMember {
		return []db.GetGroupMembersRow{}, &customerrors.UserNotAuthenticatedError{Message: "user is not a member of the group"}
	}

	users, err := repository.GetGroupMembers(ctx, groupID)
	return users, err
}

func CheckUserMemberOfGroup(ctx *gin.Context, userID uuid.UUID, groupID uuid.UUID) (bool, error) {
	isMember, err := repository.CheckUserMemberOfGroup(ctx, userID, groupID)
	return isMember, err
}

func FetchCredentialIDsWithGroupAccess(ctx *gin.Context, caller uuid.UUID, groupID uuid.UUID) ([]uuid.UUID, error) {

	credentialIDs, err := repository.FetchCredentialIDsWithGroupAccess(ctx, groupID, caller)
	return credentialIDs, err
}

func GetCredentialFieldsByGroupID(ctx *gin.Context, caller uuid.UUID, groupID uuid.UUID) ([]db.GetCredentialsFieldsByIdsRow, error) {

	isMember, err := repository.CheckUserMemberOfGroup(ctx, caller, groupID)
	if !isMember {
		return nil, &customerrors.UserNotAuthenticatedError{Message: "user does not have access to the group"}
	}
	if err != nil {
		return nil, err
	}

	credentialIDs, err := repository.FetchCredentialIDsWithGroupAccess(ctx, groupID, caller)
	if err != nil {
		return nil, err
	}

	credentialFields, err := repository.GetCredentialsFieldsByIds(ctx, credentialIDs, caller)
	if err != nil {
		return nil, err
	}

	return credentialFields, nil
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

func GetUsersOfGroups(ctx *gin.Context, groupIDs []uuid.UUID) ([]db.FetchUsersByGroupIdsRow, error) {
	users, err := repository.GetUsersOfGroups(ctx, groupIDs)
	return users, err
}

func GetCredentialGroups(ctx *gin.Context, credentialID uuid.UUID) ([]db.GetAccessTypeAndGroupsByCredentialIdRow, error) {
	groups, err := repository.GetCredentialGroups(ctx, credentialID)
	if err != nil {
		logger.Errorf(err.Error())
		return groups, err
	}
	return groups, nil
}

func GetUsersWithoutGroupAccess(ctx *gin.Context, userId uuid.UUID, groupId uuid.UUID) ([]db.GetUsersWithoutGroupAccessRow, error) {
	// Check user is can to see members of that are not in the group
	isMember, err := CheckUserMemberOfGroup(ctx, userId, groupId)
	if !isMember {
		return nil, err
	}
	users, err := repository.GetUsersWithoutGroupAccess(ctx, groupId)
	if err != nil {
		return []db.GetUsersWithoutGroupAccessRow{}, err
	}
	return users, err
}
