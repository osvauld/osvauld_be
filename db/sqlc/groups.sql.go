// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: groups.sql

package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const addGroupMember = `-- name: AddGroupMember :exec
INSERT INTO group_list (grouping_id, user_id, access_type)
VALUES ($1, $2, $3)
`

type AddGroupMemberParams struct {
	GroupingID uuid.UUID `json:"groupingId"`
	UserID     uuid.UUID `json:"userId"`
	AccessType string    `json:"accessType"`
}

func (q *Queries) AddGroupMember(ctx context.Context, arg AddGroupMemberParams) error {
	_, err := q.db.ExecContext(ctx, addGroupMember, arg.GroupingID, arg.UserID, arg.AccessType)
	return err
}

const checkUserMemberOfGroup = `-- name: CheckUserMemberOfGroup :one
SELECT EXISTS (
  SELECT 1 FROM group_list
  WHERE user_id = $1 AND grouping_id = $2
) as "exists"
`

type CheckUserMemberOfGroupParams struct {
	UserID     uuid.UUID `json:"userId"`
	GroupingID uuid.UUID `json:"groupingId"`
}

func (q *Queries) CheckUserMemberOfGroup(ctx context.Context, arg CheckUserMemberOfGroupParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, checkUserMemberOfGroup, arg.UserID, arg.GroupingID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createGroup = `-- name: CreateGroup :one
INSERT INTO groupings (name, created_by)
VALUES ($1, $2)
RETURNING id, name, created_by, created_at
`

type CreateGroupParams struct {
	Name      string    `json:"name"`
	CreatedBy uuid.UUID `json:"createdBy"`
}

type CreateGroupRow struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedBy uuid.UUID `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
}

func (q *Queries) CreateGroup(ctx context.Context, arg CreateGroupParams) (CreateGroupRow, error) {
	row := q.db.QueryRowContext(ctx, createGroup, arg.Name, arg.CreatedBy)
	var i CreateGroupRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CreatedBy,
		&i.CreatedAt,
	)
	return i, err
}

const fetchCredentialAccessTypeForGroup = `-- name: FetchCredentialAccessTypeForGroup :one
SELECT access_type FROM credential_access
WHERE group_id = $1 AND credential_id = $2
`

type FetchCredentialAccessTypeForGroupParams struct {
	GroupID      uuid.NullUUID `json:"groupId"`
	CredentialID uuid.UUID     `json:"credentialId"`
}

func (q *Queries) FetchCredentialAccessTypeForGroup(ctx context.Context, arg FetchCredentialAccessTypeForGroupParams) (string, error) {
	row := q.db.QueryRowContext(ctx, fetchCredentialAccessTypeForGroup, arg.GroupID, arg.CredentialID)
	var access_type string
	err := row.Scan(&access_type)
	return access_type, err
}

const fetchCredentialIDsWithGroupAccess = `-- name: FetchCredentialIDsWithGroupAccess :many
SELECT distinct(credential_id) from credential_access
WHERE group_id = $1 and user_id = $2
`

type FetchCredentialIDsWithGroupAccessParams struct {
	GroupID uuid.NullUUID `json:"groupId"`
	UserID  uuid.UUID     `json:"userId"`
}

func (q *Queries) FetchCredentialIDsWithGroupAccess(ctx context.Context, arg FetchCredentialIDsWithGroupAccessParams) ([]uuid.UUID, error) {
	rows, err := q.db.QueryContext(ctx, fetchCredentialIDsWithGroupAccess, arg.GroupID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []uuid.UUID{}
	for rows.Next() {
		var credential_id uuid.UUID
		if err := rows.Scan(&credential_id); err != nil {
			return nil, err
		}
		items = append(items, credential_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchGroupAccessType = `-- name: FetchGroupAccessType :one
SELECT access_type FROM group_list
WHERE user_id = $1 AND grouping_id = $2
`

type FetchGroupAccessTypeParams struct {
	UserID     uuid.UUID `json:"userId"`
	GroupingID uuid.UUID `json:"groupingId"`
}

func (q *Queries) FetchGroupAccessType(ctx context.Context, arg FetchGroupAccessTypeParams) (string, error) {
	row := q.db.QueryRowContext(ctx, fetchGroupAccessType, arg.UserID, arg.GroupingID)
	var access_type string
	err := row.Scan(&access_type)
	return access_type, err
}

const fetchUserGroups = `-- name: FetchUserGroups :many
SELECT groupings.id as "groupId", groupings.name, groupings.created_by, groupings.created_at
FROM groupings
JOIN group_list ON group_list.grouping_id = groupings.id
WHERE group_list.user_id = $1
`

type FetchUserGroupsRow struct {
	GroupId   uuid.UUID `json:"groupId"`
	Name      string    `json:"name"`
	CreatedBy uuid.UUID `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
}

func (q *Queries) FetchUserGroups(ctx context.Context, userID uuid.UUID) ([]FetchUserGroupsRow, error) {
	rows, err := q.db.QueryContext(ctx, fetchUserGroups, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []FetchUserGroupsRow{}
	for rows.Next() {
		var i FetchUserGroupsRow
		if err := rows.Scan(
			&i.GroupId,
			&i.Name,
			&i.CreatedBy,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchUsersByGroupIds = `-- name: FetchUsersByGroupIds :many
SELECT 
    g.id AS "groupId",
    json_agg(json_build_object('id', gl.user_id, 'publicKey', u.encryption_key)) AS "userDetails"
FROM 
    group_list gl
JOIN 
    groupings g ON gl.grouping_id = g.id
JOIN 
    users u ON gl.user_id = u.id
WHERE 
    g.id = ANY($1::UUID[])
GROUP BY 
    g.id
`

type FetchUsersByGroupIdsRow struct {
	GroupId     uuid.UUID       `json:"groupId"`
	UserDetails json.RawMessage `json:"userDetails"`
}

func (q *Queries) FetchUsersByGroupIds(ctx context.Context, dollar_1 []uuid.UUID) ([]FetchUsersByGroupIdsRow, error) {
	rows, err := q.db.QueryContext(ctx, fetchUsersByGroupIds, pq.Array(dollar_1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []FetchUsersByGroupIdsRow{}
	for rows.Next() {
		var i FetchUsersByGroupIdsRow
		if err := rows.Scan(&i.GroupId, &i.UserDetails); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCredentialAccessDetailsWithGroupAccess = `-- name: GetCredentialAccessDetailsWithGroupAccess :many


SELECT DISTINCT credential_id, access_type, folder_id
FROM credential_access
WHERE group_id = $1
`

type GetCredentialAccessDetailsWithGroupAccessRow struct {
	CredentialID uuid.UUID     `json:"credentialId"`
	AccessType   string        `json:"accessType"`
	FolderID     uuid.NullUUID `json:"folderId"`
}

// -----------------------------------------------------------------------------------------------------
func (q *Queries) GetCredentialAccessDetailsWithGroupAccess(ctx context.Context, groupID uuid.NullUUID) ([]GetCredentialAccessDetailsWithGroupAccessRow, error) {
	rows, err := q.db.QueryContext(ctx, getCredentialAccessDetailsWithGroupAccess, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetCredentialAccessDetailsWithGroupAccessRow{}
	for rows.Next() {
		var i GetCredentialAccessDetailsWithGroupAccessRow
		if err := rows.Scan(&i.CredentialID, &i.AccessType, &i.FolderID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFolderIDAndTypeWithGroupAccess = `-- name: GetFolderIDAndTypeWithGroupAccess :many
SELECT DISTINCT folder_id, access_type
FROM folder_access
WHERE group_id = $1
`

type GetFolderIDAndTypeWithGroupAccessRow struct {
	FolderID   uuid.UUID `json:"folderId"`
	AccessType string    `json:"accessType"`
}

func (q *Queries) GetFolderIDAndTypeWithGroupAccess(ctx context.Context, groupID uuid.NullUUID) ([]GetFolderIDAndTypeWithGroupAccessRow, error) {
	rows, err := q.db.QueryContext(ctx, getFolderIDAndTypeWithGroupAccess, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetFolderIDAndTypeWithGroupAccessRow{}
	for rows.Next() {
		var i GetFolderIDAndTypeWithGroupAccessRow
		if err := rows.Scan(&i.FolderID, &i.AccessType); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGroupMembers = `-- name: GetGroupMembers :many
SELECT users.id, users.name, users.username, COALESCE(users.encryption_key, '') as "publicKey"
FROM users
JOIN group_list ON users.id = group_list.user_id
WHERE group_list.grouping_id = $1
`

type GetGroupMembersRow struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	PublicKey string    `json:"publicKey"`
}

func (q *Queries) GetGroupMembers(ctx context.Context, groupingID uuid.UUID) ([]GetGroupMembersRow, error) {
	rows, err := q.db.QueryContext(ctx, getGroupMembers, groupingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetGroupMembersRow{}
	for rows.Next() {
		var i GetGroupMembersRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Username,
			&i.PublicKey,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGroupsWithoutAccess = `-- name: GetGroupsWithoutAccess :many
SELECT id as "groupId", name 
FROM groupings
WHERE id NOT IN (
    SELECT group_id
    FROM folder_access
    WHERE folder_id = $1 AND group_id IS NOT NULL
)
`

type GetGroupsWithoutAccessRow struct {
	GroupId uuid.UUID `json:"groupId"`
	Name    string    `json:"name"`
}

func (q *Queries) GetGroupsWithoutAccess(ctx context.Context, folderID uuid.UUID) ([]GetGroupsWithoutAccessRow, error) {
	rows, err := q.db.QueryContext(ctx, getGroupsWithoutAccess, folderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetGroupsWithoutAccessRow{}
	for rows.Next() {
		var i GetGroupsWithoutAccessRow
		if err := rows.Scan(&i.GroupId, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUsersWithoutGroupAccess = `-- name: GetUsersWithoutGroupAccess :many
SELECT u.id, u.username, u.name, COALESCE(u.encryption_key,'') as "encryptionKey"
FROM users u
LEFT JOIN group_list gl ON u.id = gl.user_id AND gl.grouping_id = $1
`

type GetUsersWithoutGroupAccessRow struct {
	ID            uuid.UUID `json:"id"`
	Username      string    `json:"username"`
	Name          string    `json:"name"`
	EncryptionKey string    `json:"encryptionKey"`
}

func (q *Queries) GetUsersWithoutGroupAccess(ctx context.Context, groupingID uuid.UUID) ([]GetUsersWithoutGroupAccessRow, error) {
	rows, err := q.db.QueryContext(ctx, getUsersWithoutGroupAccess, groupingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetUsersWithoutGroupAccessRow{}
	for rows.Next() {
		var i GetUsersWithoutGroupAccessRow
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Name,
			&i.EncryptionKey,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
