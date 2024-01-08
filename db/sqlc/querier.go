// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package db

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
)

type Querier interface {
	AddCredential(ctx context.Context, dollar_1 json.RawMessage) (interface{}, error)
	AddFolderAccess(ctx context.Context, arg AddFolderAccessParams) error
	AddGroupMemberRecord(ctx context.Context, arg AddGroupMemberRecordParams) error
	AddToAccessList(ctx context.Context, arg AddToAccessListParams) (uuid.UUID, error)
	CheckTempPassword(ctx context.Context, arg CheckTempPasswordParams) (int64, error)
	CheckUserMemberOfGroup(ctx context.Context, arg CheckUserMemberOfGroupParams) (bool, error)
	CreateChallenge(ctx context.Context, arg CreateChallengeParams) (SessionTable, error)
	// sql/create_credential.sql
	CreateCredential(ctx context.Context, arg CreateCredentialParams) (uuid.UUID, error)
	CreateEncryptedData(ctx context.Context, arg CreateEncryptedDataParams) (uuid.UUID, error)
	CreateFolder(ctx context.Context, arg CreateFolderParams) (uuid.UUID, error)
	CreateGroup(ctx context.Context, arg CreateGroupParams) (uuid.UUID, error)
	CreateUnencryptedData(ctx context.Context, arg CreateUnencryptedDataParams) (uuid.UUID, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (uuid.UUID, error)
	FetchAccessibleAndCreatedFoldersByUser(ctx context.Context, createdBy uuid.UUID) ([]FetchAccessibleAndCreatedFoldersByUserRow, error)
	FetchChallenge(ctx context.Context, userID uuid.UUID) (string, error)
	FetchCredentialAccessTypeForGroupMember(ctx context.Context, arg FetchCredentialAccessTypeForGroupMemberParams) (string, error)
	FetchCredentialDataByID(ctx context.Context, id uuid.UUID) (Credential, error)
	FetchCredentialIDsWithGroupAccess(ctx context.Context, groupID uuid.NullUUID) ([]uuid.UUID, error)
	FetchCredentialsByUserAndFolder(ctx context.Context, arg FetchCredentialsByUserAndFolderParams) ([]FetchCredentialsByUserAndFolderRow, error)
	FetchEncryptedFieldsByCredentialIDAndUserID(ctx context.Context, arg FetchEncryptedFieldsByCredentialIDAndUserIDParams) ([]FetchEncryptedFieldsByCredentialIDAndUserIDRow, error)
	FetchGroupAccessType(ctx context.Context, arg FetchGroupAccessTypeParams) (string, error)
	FetchUnencryptedFieldsByCredentialID(ctx context.Context, credentialID uuid.UUID) ([]FetchUnencryptedFieldsByCredentialIDRow, error)
	FetchUserGroups(ctx context.Context, userID uuid.UUID) ([]FetchUserGroupsRow, error)
	GetAllUsers(ctx context.Context) ([]GetAllUsersRow, error)
	GetCredentialAccessForUser(ctx context.Context, arg GetCredentialAccessForUserParams) ([]GetCredentialAccessForUserRow, error)
	GetCredentialDetails(ctx context.Context, id uuid.UUID) (GetCredentialDetailsRow, error)
	GetCredentialIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	GetCredentialUnencryptedData(ctx context.Context, credentialID uuid.UUID) ([]GetCredentialUnencryptedDataRow, error)
	GetEncryptedCredentialsByFolder(ctx context.Context, arg GetEncryptedCredentialsByFolderParams) ([]GetEncryptedCredentialsByFolderRow, error)
	GetEncryptedDataByCredentialIds(ctx context.Context, arg GetEncryptedDataByCredentialIdsParams) ([]GetEncryptedDataByCredentialIdsRow, error)
	GetGroupMembers(ctx context.Context, groupingID uuid.UUID) ([]GetGroupMembersRow, error)
	GetSharedUsers(ctx context.Context, folderID uuid.UUID) ([]GetSharedUsersRow, error)
	GetUserByPublicKey(ctx context.Context, eccPubKey sql.NullString) (uuid.UUID, error)
	GetUserByUsername(ctx context.Context, username string) (GetUserByUsernameRow, error)
	GetUserEncryptedData(ctx context.Context, arg GetUserEncryptedDataParams) ([]GetUserEncryptedDataRow, error)
	GetUsersByCredential(ctx context.Context, credentialID uuid.UUID) ([]GetUsersByCredentialRow, error)
	GetUsersByFolder(ctx context.Context, folderID uuid.UUID) ([]GetUsersByFolderRow, error)
	IsFolderOwner(ctx context.Context, arg IsFolderOwnerParams) (bool, error)
	ShareSecret(ctx context.Context, dollar_1 json.RawMessage) error
	UpdateKeys(ctx context.Context, arg UpdateKeysParams) error
}

var _ Querier = (*Queries)(nil)
