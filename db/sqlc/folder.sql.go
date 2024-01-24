// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: folder.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const addFolderAccess = `-- name: AddFolderAccess :exec
INSERT INTO folder_access (folder_id, user_id, access_type)
VALUES ($1, $2, $3)
`

type AddFolderAccessParams struct {
	FolderID   uuid.UUID `json:"folder_id"`
	UserID     uuid.UUID `json:"user_id"`
	AccessType string    `json:"access_type"`
}

func (q *Queries) AddFolderAccess(ctx context.Context, arg AddFolderAccessParams) error {
	_, err := q.db.ExecContext(ctx, addFolderAccess, arg.FolderID, arg.UserID, arg.AccessType)
	return err
}

const addFolderAccessWithGroup = `-- name: AddFolderAccessWithGroup :exec
INSERT INTO folder_access (folder_id, user_id, access_type, group_id)
VALUES ($1, $2, $3, $4)
`

type AddFolderAccessWithGroupParams struct {
	FolderID   uuid.UUID     `json:"folder_id"`
	UserID     uuid.UUID     `json:"user_id"`
	AccessType string        `json:"access_type"`
	GroupID    uuid.NullUUID `json:"group_id"`
}

func (q *Queries) AddFolderAccessWithGroup(ctx context.Context, arg AddFolderAccessWithGroupParams) error {
	_, err := q.db.ExecContext(ctx, addFolderAccessWithGroup,
		arg.FolderID,
		arg.UserID,
		arg.AccessType,
		arg.GroupID,
	)
	return err
}

const createFolder = `-- name: CreateFolder :one
WITH new_folder AS (
  INSERT INTO folders (name, description, created_by)
  VALUES ($1, $2, $3)
  RETURNING id
),
folder_access_insert AS (
  INSERT INTO folder_access (folder_id, user_id, access_type)
  SELECT id, $3, 'owner' FROM new_folder
)
SELECT id FROM new_folder
`

type CreateFolderParams struct {
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	CreatedBy   uuid.UUID      `json:"created_by"`
}

func (q *Queries) CreateFolder(ctx context.Context, arg CreateFolderParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createFolder, arg.Name, arg.Description, arg.CreatedBy)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const fetchAccessibleAndCreatedFoldersByUser = `-- name: FetchAccessibleAndCreatedFoldersByUser :many
WITH unique_credential_ids AS (
  SELECT DISTINCT credential_id
  FROM access_list
  WHERE user_id = $1
),
unique_folder_ids AS (
  SELECT DISTINCT folder_id
  FROM credentials
  WHERE id IN (SELECT credential_id FROM unique_credential_ids)
)
SELECT 
    id, 
    name, 
    COALESCE(description, '') AS description 
FROM folders f
WHERE f.id IN (SELECT folder_id FROM unique_folder_ids)
   OR f.created_by = $1
`

type FetchAccessibleAndCreatedFoldersByUserRow struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func (q *Queries) FetchAccessibleAndCreatedFoldersByUser(ctx context.Context, createdBy uuid.UUID) ([]FetchAccessibleAndCreatedFoldersByUserRow, error) {
	rows, err := q.db.QueryContext(ctx, fetchAccessibleAndCreatedFoldersByUser, createdBy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []FetchAccessibleAndCreatedFoldersByUserRow{}
	for rows.Next() {
		var i FetchAccessibleAndCreatedFoldersByUserRow
		if err := rows.Scan(&i.ID, &i.Name, &i.Description); err != nil {
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

const getAccessTypeAndUserByFolder = `-- name: GetAccessTypeAndUserByFolder :many
SELECT user_id, access_type
FROM folder_access
WHERE folder_id = $1
`

type GetAccessTypeAndUserByFolderRow struct {
	UserID     uuid.UUID `json:"user_id"`
	AccessType string    `json:"access_type"`
}

func (q *Queries) GetAccessTypeAndUserByFolder(ctx context.Context, folderID uuid.UUID) ([]GetAccessTypeAndUserByFolderRow, error) {
	rows, err := q.db.QueryContext(ctx, getAccessTypeAndUserByFolder, folderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAccessTypeAndUserByFolderRow{}
	for rows.Next() {
		var i GetAccessTypeAndUserByFolderRow
		if err := rows.Scan(&i.UserID, &i.AccessType); err != nil {
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

const getFolderAccessForUser = `-- name: GetFolderAccessForUser :many
SELECT access_type FROM folder_access
WHERE folder_id = $1 AND user_id = $2
`

type GetFolderAccessForUserParams struct {
	FolderID uuid.UUID `json:"folder_id"`
	UserID   uuid.UUID `json:"user_id"`
}

func (q *Queries) GetFolderAccessForUser(ctx context.Context, arg GetFolderAccessForUserParams) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getFolderAccessForUser, arg.FolderID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var access_type string
		if err := rows.Scan(&access_type); err != nil {
			return nil, err
		}
		items = append(items, access_type)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSharedUsers = `-- name: GetSharedUsers :many
SELECT users.id, users.name, users.username, COALESCE(users.rsa_pub_key,'') as "publicKey", folder_access.access_type as "accessType"
FROM folder_access
JOIN users ON folder_access.user_id = users.id
WHERE folder_access.folder_id = $1
`

type GetSharedUsersRow struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Username   string    `json:"username"`
	PublicKey  string    `json:"publicKey"`
	AccessType string    `json:"accessType"`
}

func (q *Queries) GetSharedUsers(ctx context.Context, folderID uuid.UUID) ([]GetSharedUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, getSharedUsers, folderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetSharedUsersRow{}
	for rows.Next() {
		var i GetSharedUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Username,
			&i.PublicKey,
			&i.AccessType,
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

const isFolderOwner = `-- name: IsFolderOwner :one
SELECT EXISTS (
  SELECT 1 FROM folder_access
  WHERE folder_id = $1 AND user_id = $2 AND access_type = 'owner'
)
`

type IsFolderOwnerParams struct {
	FolderID uuid.UUID `json:"folder_id"`
	UserID   uuid.UUID `json:"user_id"`
}

func (q *Queries) IsFolderOwner(ctx context.Context, arg IsFolderOwnerParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, isFolderOwner, arg.FolderID, arg.UserID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const isUserManagerOrOwner = `-- name: IsUserManagerOrOwner :one
SELECT EXISTS (
  SELECT 1 FROM folder_access
  WHERE folder_id = $1 AND user_id = $2 AND access_type IN ('owner', 'manager')
)
`

type IsUserManagerOrOwnerParams struct {
	FolderID uuid.UUID `json:"folder_id"`
	UserID   uuid.UUID `json:"user_id"`
}

func (q *Queries) IsUserManagerOrOwner(ctx context.Context, arg IsUserManagerOrOwnerParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, isUserManagerOrOwner, arg.FolderID, arg.UserID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}
