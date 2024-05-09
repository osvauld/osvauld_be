// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Querier interface {
	AddCredentialAccess(ctx context.Context, arg AddCredentialAccessParams) (uuid.UUID, error)
	AddEnvironment(ctx context.Context, arg AddEnvironmentParams) (uuid.UUID, error)
	AddField(ctx context.Context, arg AddFieldParams) (uuid.UUID, error)
	AddFolder(ctx context.Context, arg AddFolderParams) (AddFolderRow, error)
	AddFolderAccess(ctx context.Context, arg AddFolderAccessParams) error
	AddGroupMember(ctx context.Context, arg AddGroupMemberParams) error
	CheckCredentialAccessEntryExists(ctx context.Context, arg CheckCredentialAccessEntryExistsParams) (bool, error)
	CheckCredentialExistsForEnv(ctx context.Context, arg CheckCredentialExistsForEnvParams) (bool, error)
	CheckFieldEntryExists(ctx context.Context, arg CheckFieldEntryExistsParams) (bool, error)
	CheckFolderAccessEntryExists(ctx context.Context, arg CheckFolderAccessEntryExistsParams) (bool, error)
	CheckIfUsersExist(ctx context.Context) (bool, error)
	CheckNameExist(ctx context.Context, name string) (bool, error)
	CheckUserAdminOfGroup(ctx context.Context, arg CheckUserAdminOfGroupParams) (bool, error)
	CheckUserMemberOfGroup(ctx context.Context, arg CheckUserMemberOfGroupParams) (bool, error)
	CheckUsernameExist(ctx context.Context, username string) (bool, error)
	CreateChallenge(ctx context.Context, arg CreateChallengeParams) (SessionTable, error)
	CreateCliUser(ctx context.Context, arg CreateCliUserParams) (uuid.UUID, error)
	// sql/create_credential.sql
	CreateCredential(ctx context.Context, arg CreateCredentialParams) (uuid.UUID, error)
	CreateEnvFields(ctx context.Context, arg CreateEnvFieldsParams) (uuid.UUID, error)
	CreateGroup(ctx context.Context, arg CreateGroupParams) (CreateGroupRow, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (uuid.UUID, error)
	DeleteAccessRemovedFields(ctx context.Context) error
	DeleteCredentialFields(ctx context.Context, credentialID uuid.UUID) error
	EditCredentialAccessForGroup(ctx context.Context, arg EditCredentialAccessForGroupParams) error
	EditCredentialAccessForGroupWithFolderID(ctx context.Context, arg EditCredentialAccessForGroupWithFolderIDParams) error
	EditCredentialAccessForUser(ctx context.Context, arg EditCredentialAccessForUserParams) error
	EditCredentialAccessForUserWithFolderID(ctx context.Context, arg EditCredentialAccessForUserWithFolderIDParams) error
	EditCredentialDetails(ctx context.Context, arg EditCredentialDetailsParams) error
	EditFolder(ctx context.Context, arg EditFolderParams) error
	EditFolderAccessForGroup(ctx context.Context, arg EditFolderAccessForGroupParams) error
	EditFolderAccessForUser(ctx context.Context, arg EditFolderAccessForUserParams) error
	EditGroup(ctx context.Context, arg EditGroupParams) error
	FetchAccessibleFoldersForUser(ctx context.Context, userID uuid.UUID) ([]FetchAccessibleFoldersForUserRow, error)
	FetchChallenge(ctx context.Context, userID uuid.UUID) (string, error)
	FetchCredentialAccessTypeForGroup(ctx context.Context, arg FetchCredentialAccessTypeForGroupParams) (string, error)
	FetchCredentialDetailsForUserByFolderId(ctx context.Context, arg FetchCredentialDetailsForUserByFolderIdParams) ([]FetchCredentialDetailsForUserByFolderIdRow, error)
	FetchCredentialIDsWithGroupAccess(ctx context.Context, arg FetchCredentialIDsWithGroupAccessParams) ([]uuid.UUID, error)
	FetchUserGroups(ctx context.Context, userID uuid.UUID) ([]FetchUserGroupsRow, error)
	FetchUsersByGroupIds(ctx context.Context, dollar_1 []uuid.UUID) ([]FetchUsersByGroupIdsRow, error)
	GetAccessTypeAndGroupsByCredentialId(ctx context.Context, credentialID uuid.UUID) ([]GetAccessTypeAndGroupsByCredentialIdRow, error)
	GetAccessTypeAndUserByFolder(ctx context.Context, folderID uuid.UUID) ([]GetAccessTypeAndUserByFolderRow, error)
	GetAllFieldsForCredentialIDs(ctx context.Context, arg GetAllFieldsForCredentialIDsParams) ([]GetAllFieldsForCredentialIDsRow, error)
	GetAllSignedUpUsers(ctx context.Context) ([]GetAllSignedUpUsersRow, error)
	GetAllUrlsForUser(ctx context.Context, userID uuid.UUID) ([]GetAllUrlsForUserRow, error)
	GetAllUsers(ctx context.Context) ([]GetAllUsersRow, error)
	//-----------------------------------------------------------------------------------------------------
	GetCredentialAccessDetailsWithGroupAccess(ctx context.Context, groupID uuid.NullUUID) ([]GetCredentialAccessDetailsWithGroupAccessRow, error)
	GetCredentialAccessTypeForUser(ctx context.Context, arg GetCredentialAccessTypeForUserParams) ([]GetCredentialAccessTypeForUserRow, error)
	GetCredentialDataByID(ctx context.Context, id uuid.UUID) (GetCredentialDataByIDRow, error)
	GetCredentialDetailsByIDs(ctx context.Context, credentialids []uuid.UUID) ([]GetCredentialDetailsByIDsRow, error)
	GetCredentialGroups(ctx context.Context, credentialID uuid.UUID) ([]GetCredentialGroupsRow, error)
	GetCredentialIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	GetCredentialIdsByFolder(ctx context.Context, arg GetCredentialIdsByFolderParams) ([]uuid.UUID, error)
	GetCredentialUsersForDataSync(ctx context.Context, credentialID uuid.UUID) ([]GetCredentialUsersForDataSyncRow, error)
	GetCredentialUsersWithDirectAccess(ctx context.Context, credentialID uuid.UUID) ([]GetCredentialUsersWithDirectAccessRow, error)
	GetCredentialsForSearchByUserID(ctx context.Context, userID uuid.UUID) ([]GetCredentialsForSearchByUserIDRow, error)
	GetFolderAccessForUser(ctx context.Context, arg GetFolderAccessForUserParams) ([]string, error)
	GetFolderGroups(ctx context.Context, folderID uuid.UUID) ([]GetFolderGroupsRow, error)
	GetFolderIDAndTypeWithGroupAccess(ctx context.Context, groupID uuid.NullUUID) ([]GetFolderIDAndTypeWithGroupAccessRow, error)
	GetFolderUsersForDataSync(ctx context.Context, folderID uuid.UUID) ([]GetFolderUsersForDataSyncRow, error)
	GetFolderUsersWithDirectAccess(ctx context.Context, folderID uuid.UUID) ([]GetFolderUsersWithDirectAccessRow, error)
	GetGroupMembers(ctx context.Context, groupingID uuid.UUID) ([]GetGroupMembersRow, error)
	GetGroupsWithoutAccess(ctx context.Context, arg GetGroupsWithoutAccessParams) ([]GetGroupsWithoutAccessRow, error)
	GetNonSensitiveFieldsForCredentialIDs(ctx context.Context, arg GetNonSensitiveFieldsForCredentialIDsParams) ([]GetNonSensitiveFieldsForCredentialIDsRow, error)
	GetRegistrationChallenge(ctx context.Context, username string) (GetRegistrationChallengeRow, error)
	GetSensitiveFields(ctx context.Context, arg GetSensitiveFieldsParams) ([]GetSensitiveFieldsRow, error)
	GetSharedGroupsForFolder(ctx context.Context, folderID uuid.UUID) ([]GetSharedGroupsForFolderRow, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (GetUserByIDRow, error)
	GetUserByPublicKey(ctx context.Context, deviceKey sql.NullString) (uuid.UUID, error)
	GetUserByUsername(ctx context.Context, username string) (GetUserByUsernameRow, error)
	GetUserDeviceKey(ctx context.Context, id uuid.UUID) (string, error)
	GetUserTempPassword(ctx context.Context, username string) (GetUserTempPasswordRow, error)
	GetUserType(ctx context.Context, id uuid.UUID) (string, error)
	GetUsersByCredential(ctx context.Context, credentialID uuid.UUID) ([]GetUsersByCredentialRow, error)
	GetUsersWithoutGroupAccess(ctx context.Context, groupingID uuid.UUID) ([]GetUsersWithoutGroupAccessRow, error)
	HasManageAccessForCredential(ctx context.Context, arg HasManageAccessForCredentialParams) (bool, error)
	HasManageAccessForFolder(ctx context.Context, arg HasManageAccessForFolderParams) (bool, error)
	HasReadAccessForCredential(ctx context.Context, arg HasReadAccessForCredentialParams) (bool, error)
	HasReadAccessForFolder(ctx context.Context, arg HasReadAccessForFolderParams) (bool, error)
	IsUserManagerOrOwner(ctx context.Context, arg IsUserManagerOrOwnerParams) (bool, error)
	RemoveCredential(ctx context.Context, id uuid.UUID) error
	RemoveCredentialAccessForGroupMember(ctx context.Context, arg RemoveCredentialAccessForGroupMemberParams) error
	RemoveCredentialAccessForGroups(ctx context.Context, arg RemoveCredentialAccessForGroupsParams) error
	RemoveCredentialAccessForGroupsWithFolderID(ctx context.Context, arg RemoveCredentialAccessForGroupsWithFolderIDParams) error
	RemoveCredentialAccessForUsers(ctx context.Context, arg RemoveCredentialAccessForUsersParams) error
	RemoveCredentialAccessForUsersWithFolderID(ctx context.Context, arg RemoveCredentialAccessForUsersWithFolderIDParams) error
	RemoveCredentialFieldsForUsers(ctx context.Context, arg RemoveCredentialFieldsForUsersParams) error
	RemoveFolder(ctx context.Context, id uuid.UUID) error
	RemoveFolderAccessForGroupMember(ctx context.Context, arg RemoveFolderAccessForGroupMemberParams) error
	RemoveFolderAccessForGroups(ctx context.Context, arg RemoveFolderAccessForGroupsParams) error
	RemoveFolderAccessForUsers(ctx context.Context, arg RemoveFolderAccessForUsersParams) error
	RemoveGroup(ctx context.Context, id uuid.UUID) error
	RemoveUserFromGroupList(ctx context.Context, arg RemoveUserFromGroupListParams) error
	RemoveUserFromOrg(ctx context.Context, id uuid.UUID) error
	UpdateKeys(ctx context.Context, arg UpdateKeysParams) error
	UpdateRegistrationChallenge(ctx context.Context, arg UpdateRegistrationChallengeParams) error
}

var _ Querier = (*Queries)(nil)
