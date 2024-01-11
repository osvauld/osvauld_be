// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: groups.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const addGroupMemberRecord = `-- name: AddGroupMemberRecord :exec
INSERT INTO group_list (grouping_id, user_id, access_type)
VALUES ($1, $2, $3)
`

type AddGroupMemberRecordParams struct {
	GroupingID uuid.UUID `json:"grouping_id"`
	UserID     uuid.UUID `json:"user_id"`
	AccessType string    `json:"access_type"`
}

func (q *Queries) AddGroupMemberRecord(ctx context.Context, arg AddGroupMemberRecordParams) error {
	_, err := q.db.ExecContext(ctx, addGroupMemberRecord, arg.GroupingID, arg.UserID, arg.AccessType)
	return err
}

const checkUserMemberOfGroup = `-- name: CheckUserMemberOfGroup :one
SELECT EXISTS (
  SELECT 1 FROM group_list
  WHERE user_id = $1 AND grouping_id = $2
) as "exists"
`

type CheckUserMemberOfGroupParams struct {
	UserID     uuid.UUID `json:"user_id"`
	GroupingID uuid.UUID `json:"grouping_id"`
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
RETURNING id
`

type CreateGroupParams struct {
	Name      string    `json:"name"`
	CreatedBy uuid.UUID `json:"created_by"`
}

func (q *Queries) CreateGroup(ctx context.Context, arg CreateGroupParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createGroup, arg.Name, arg.CreatedBy)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const fetchCredentialAccessTypeForGroupMember = `-- name: FetchCredentialAccessTypeForGroupMember :one
SELECT access_type FROM access_list
WHERE group_id = $1 AND credential_id = $2 AND user_id = $3
`

type FetchCredentialAccessTypeForGroupMemberParams struct {
	GroupID      uuid.NullUUID `json:"group_id"`
	CredentialID uuid.UUID     `json:"credential_id"`
	UserID       uuid.UUID     `json:"user_id"`
}

func (q *Queries) FetchCredentialAccessTypeForGroupMember(ctx context.Context, arg FetchCredentialAccessTypeForGroupMemberParams) (string, error) {
	row := q.db.QueryRowContext(ctx, fetchCredentialAccessTypeForGroupMember, arg.GroupID, arg.CredentialID, arg.UserID)
	var access_type string
	err := row.Scan(&access_type)
	return access_type, err
}

const fetchCredentialIDsWithGroupAccess = `-- name: FetchCredentialIDsWithGroupAccess :many
SELECT distinct(credential_id) from access_list
WHERE group_id = $1
`

func (q *Queries) FetchCredentialIDsWithGroupAccess(ctx context.Context, groupID uuid.NullUUID) ([]uuid.UUID, error) {
	rows, err := q.db.QueryContext(ctx, fetchCredentialIDsWithGroupAccess, groupID)
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
	UserID     uuid.UUID `json:"user_id"`
	GroupingID uuid.UUID `json:"grouping_id"`
}

func (q *Queries) FetchGroupAccessType(ctx context.Context, arg FetchGroupAccessTypeParams) (string, error) {
	row := q.db.QueryRowContext(ctx, fetchGroupAccessType, arg.UserID, arg.GroupingID)
	var access_type string
	err := row.Scan(&access_type)
	return access_type, err
}

const fetchUserGroups = `-- name: FetchUserGroups :many
SELECT groupings.id, groupings.name, groupings.created_by, groupings.created_at
FROM groupings
JOIN group_list ON group_list.grouping_id = groupings.id
WHERE group_list.user_id = $1
`

type FetchUserGroupsRow struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedBy uuid.UUID `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
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
			&i.ID,
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

const getGroupMembers = `-- name: GetGroupMembers :many
SELECT users.id, users.name, users.username, COALESCE(users.rsa_pub_key, '') as "publicKey"
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
