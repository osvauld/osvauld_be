// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: credential.sql

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const createCredential = `-- name: CreateCredential :one
INSERT INTO
    credentials (NAME, description, credential_type, folder_id, created_by, domain)
VALUES
    ($1, $2, $3, $4, $5, $6) RETURNING id
`

type CreateCredentialParams struct {
	Name           string         `json:"name"`
	Description    sql.NullString `json:"description"`
	CredentialType string         `json:"credentialType"`
	FolderID       uuid.UUID      `json:"folderId"`
	CreatedBy      uuid.NullUUID  `json:"createdBy"`
	Domain         sql.NullString `json:"domain"`
}

// sql/create_credential.sql
func (q *Queries) CreateCredential(ctx context.Context, arg CreateCredentialParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createCredential,
		arg.Name,
		arg.Description,
		arg.CredentialType,
		arg.FolderID,
		arg.CreatedBy,
		arg.Domain,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const editCredentialDetails = `-- name: EditCredentialDetails :exec
UPDATE
    credentials
SET
    name = $2,
    description = $3,
    credential_type = $4,
    updated_at = NOW(),
    updated_by = $5,
    domain = $6
WHERE
    id = $1
`

type EditCredentialDetailsParams struct {
	ID             uuid.UUID      `json:"id"`
	Name           string         `json:"name"`
	Description    sql.NullString `json:"description"`
	CredentialType string         `json:"credentialType"`
	UpdatedBy      uuid.NullUUID  `json:"updatedBy"`
	Domain         sql.NullString `json:"domain"`
}

func (q *Queries) EditCredentialDetails(ctx context.Context, arg EditCredentialDetailsParams) error {
	_, err := q.db.ExecContext(ctx, editCredentialDetails,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.CredentialType,
		arg.UpdatedBy,
		arg.Domain,
	)
	return err
}

const fetchCredentialDetailsForUserByFolderId = `-- name: FetchCredentialDetailsForUserByFolderId :many
SELECT
    C.id AS "credentialID",
    C.name,
    COALESCE(C.description, '') AS "description",
    C.credential_type AS "credentialType",
    C.created_at AS "createdAt",
    C.updated_at AS "updatedAt",
    C.created_by AS "createdBy",
    C.updated_by AS "updatedBy",
    A.access_type AS "accessType"
FROM
    credentials AS C,
    credential_access AS A
WHERE
    C.id = A .credential_id
    AND C.folder_id = $1
    AND A.user_id = $2
`

type FetchCredentialDetailsForUserByFolderIdParams struct {
	FolderID uuid.UUID `json:"folderId"`
	UserID   uuid.UUID `json:"userId"`
}

type FetchCredentialDetailsForUserByFolderIdRow struct {
	CredentialID   uuid.UUID     `json:"credentialID"`
	Name           string        `json:"name"`
	Description    string        `json:"description"`
	CredentialType string        `json:"credentialType"`
	CreatedAt      time.Time     `json:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt"`
	CreatedBy      uuid.NullUUID `json:"createdBy"`
	UpdatedBy      uuid.NullUUID `json:"updatedBy"`
	AccessType     string        `json:"accessType"`
}

func (q *Queries) FetchCredentialDetailsForUserByFolderId(ctx context.Context, arg FetchCredentialDetailsForUserByFolderIdParams) ([]FetchCredentialDetailsForUserByFolderIdRow, error) {
	rows, err := q.db.QueryContext(ctx, fetchCredentialDetailsForUserByFolderId, arg.FolderID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []FetchCredentialDetailsForUserByFolderIdRow{}
	for rows.Next() {
		var i FetchCredentialDetailsForUserByFolderIdRow
		if err := rows.Scan(
			&i.CredentialID,
			&i.Name,
			&i.Description,
			&i.CredentialType,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.CreatedBy,
			&i.UpdatedBy,
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

const getAccessTypeAndGroupsByCredentialId = `-- name: GetAccessTypeAndGroupsByCredentialId :many
    SELECT DISTINCT
        al.group_id, 
        g.name,
        al.access_type
    FROM 
        credential_access al
    JOIN 
        groupings g ON al.group_id = g.id
    WHERE 
        al.credential_id = $1
`

type GetAccessTypeAndGroupsByCredentialIdRow struct {
	GroupID    uuid.NullUUID `json:"groupId"`
	Name       string        `json:"name"`
	AccessType string        `json:"accessType"`
}

func (q *Queries) GetAccessTypeAndGroupsByCredentialId(ctx context.Context, credentialID uuid.UUID) ([]GetAccessTypeAndGroupsByCredentialIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getAccessTypeAndGroupsByCredentialId, credentialID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAccessTypeAndGroupsByCredentialIdRow{}
	for rows.Next() {
		var i GetAccessTypeAndGroupsByCredentialIdRow
		if err := rows.Scan(&i.GroupID, &i.Name, &i.AccessType); err != nil {
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

const getAllUrlsForUser = `-- name: GetAllUrlsForUser :many
SELECT DISTINCT
    fv.field_value AS value,
    fd.credential_id AS "credentialId"
FROM 
    field_values fv
JOIN 
    field_data fd ON fv.field_id = fd.id
WHERE 
    fv.user_id = $1 AND fd.field_name = 'Domain'
`

type GetAllUrlsForUserRow struct {
	Value        string    `json:"value"`
	CredentialId uuid.UUID `json:"credentialId"`
}

func (q *Queries) GetAllUrlsForUser(ctx context.Context, userID uuid.UUID) ([]GetAllUrlsForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllUrlsForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllUrlsForUserRow{}
	for rows.Next() {
		var i GetAllUrlsForUserRow
		if err := rows.Scan(&i.Value, &i.CredentialId); err != nil {
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

const getCredentialDataByID = `-- name: GetCredentialDataByID :one
SELECT
    id,
    name,
    description,
    folder_id,
    credential_type,
    created_by,
    created_at,
    updated_at,
    updated_by
FROM
    credentials
WHERE
    id = $1
`

type GetCredentialDataByIDRow struct {
	ID             uuid.UUID      `json:"id"`
	Name           string         `json:"name"`
	Description    sql.NullString `json:"description"`
	FolderID       uuid.UUID      `json:"folderId"`
	CredentialType string         `json:"credentialType"`
	CreatedBy      uuid.NullUUID  `json:"createdBy"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	UpdatedBy      uuid.NullUUID  `json:"updatedBy"`
}

func (q *Queries) GetCredentialDataByID(ctx context.Context, id uuid.UUID) (GetCredentialDataByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getCredentialDataByID, id)
	var i GetCredentialDataByIDRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.FolderID,
		&i.CredentialType,
		&i.CreatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UpdatedBy,
	)
	return i, err
}

const getCredentialDetailsByIDs = `-- name: GetCredentialDetailsByIDs :many
SELECT
    id,
    name,
    description,
    folder_id,
    credential_type,
    created_by,
    created_at,
    updated_at,
    updated_by
FROM
    credentials
WHERE
    id = ANY($1::UUID[])
`

type GetCredentialDetailsByIDsRow struct {
	ID             uuid.UUID      `json:"id"`
	Name           string         `json:"name"`
	Description    sql.NullString `json:"description"`
	FolderID       uuid.UUID      `json:"folderId"`
	CredentialType string         `json:"credentialType"`
	CreatedBy      uuid.NullUUID  `json:"createdBy"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	UpdatedBy      uuid.NullUUID  `json:"updatedBy"`
}

func (q *Queries) GetCredentialDetailsByIDs(ctx context.Context, credentialids []uuid.UUID) ([]GetCredentialDetailsByIDsRow, error) {
	rows, err := q.db.QueryContext(ctx, getCredentialDetailsByIDs, pq.Array(credentialids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetCredentialDetailsByIDsRow{}
	for rows.Next() {
		var i GetCredentialDetailsByIDsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.FolderID,
			&i.CredentialType,
			&i.CreatedBy,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UpdatedBy,
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

const getCredentialIdsByFolder = `-- name: GetCredentialIdsByFolder :many
SELECT 
    DISTINCT c.id AS "credentialId"
FROM 
    credentials c
JOIN 
    credential_access a ON c.id = a.credential_id
WHERE 
    a.user_id = $1
    AND c.folder_id = $2
`

type GetCredentialIdsByFolderParams struct {
	UserID   uuid.UUID `json:"userId"`
	FolderID uuid.UUID `json:"folderId"`
}

func (q *Queries) GetCredentialIdsByFolder(ctx context.Context, arg GetCredentialIdsByFolderParams) ([]uuid.UUID, error) {
	rows, err := q.db.QueryContext(ctx, getCredentialIdsByFolder, arg.UserID, arg.FolderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []uuid.UUID{}
	for rows.Next() {
		var credentialId uuid.UUID
		if err := rows.Scan(&credentialId); err != nil {
			return nil, err
		}
		items = append(items, credentialId)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCredentialsForSearchByUserID = `-- name: GetCredentialsForSearchByUserID :many
SELECT DISTINCT
    c.id as "credentialId", 
    c.name, 
    COALESCE(c.description, '') AS description,
    COALESCE(c.domain, '') AS domain,
    c.folder_id,
    COALESCE(f.type, '' ) AS "folderType",
    COALESCE(f.name, '') AS folder_name
FROM 
    credentials c
JOIN 
    credential_access ca ON c.id = ca.credential_id
LEFT JOIN 
    folders f ON c.folder_id = f.id
WHERE 
    ca.user_id = $1
`

type GetCredentialsForSearchByUserIDRow struct {
	CredentialId uuid.UUID `json:"credentialId"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Domain       string    `json:"domain"`
	FolderID     uuid.UUID `json:"folderId"`
	FolderType   string    `json:"folderType"`
	FolderName   string    `json:"folderName"`
}

func (q *Queries) GetCredentialsForSearchByUserID(ctx context.Context, userID uuid.UUID) ([]GetCredentialsForSearchByUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getCredentialsForSearchByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetCredentialsForSearchByUserIDRow{}
	for rows.Next() {
		var i GetCredentialsForSearchByUserIDRow
		if err := rows.Scan(
			&i.CredentialId,
			&i.Name,
			&i.Description,
			&i.Domain,
			&i.FolderID,
			&i.FolderType,
			&i.FolderName,
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

const removeCredential = `-- name: RemoveCredential :exec
DELETE FROM 
    credentials
WHERE 
    id = $1
`

func (q *Queries) RemoveCredential(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, removeCredential, id)
	return err
}
