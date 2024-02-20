// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: credential.sql

package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const createCredential = `-- name: CreateCredential :one
INSERT INTO
    credentials (NAME, description, credential_type, folder_id, created_by)
VALUES
    ($1, $2, $3, $4, $5) RETURNING id
`

type CreateCredentialParams struct {
	Name           string         `json:"name"`
	Description    sql.NullString `json:"description"`
	CredentialType string         `json:"credentialType"`
	FolderID       uuid.UUID      `json:"folderId"`
	CreatedBy      uuid.UUID      `json:"createdBy"`
}

// sql/create_credential.sql
func (q *Queries) CreateCredential(ctx context.Context, arg CreateCredentialParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createCredential,
		arg.Name,
		arg.Description,
		arg.CredentialType,
		arg.FolderID,
		arg.CreatedBy,
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
    credential_type = $4
WHERE
    id = $1
`

type EditCredentialDetailsParams struct {
	ID             uuid.UUID      `json:"id"`
	Name           string         `json:"name"`
	Description    sql.NullString `json:"description"`
	CredentialType string         `json:"credentialType"`
}

func (q *Queries) EditCredentialDetails(ctx context.Context, arg EditCredentialDetailsParams) error {
	_, err := q.db.ExecContext(ctx, editCredentialDetails,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.CredentialType,
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
	CredentialID   uuid.UUID `json:"credentialID"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	CredentialType string    `json:"credentialType"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	CreatedBy      uuid.UUID `json:"createdBy"`
	AccessType     string    `json:"accessType"`
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

const getAccessTypeAndUsersByCredentialId = `-- name: GetAccessTypeAndUsersByCredentialId :many
SELECT 
    al.user_id as "id",
    u.name, 
    al.access_type,
    COALESCE(u.rsa_pub_key, '') AS "publicKey"
FROM 
    credential_access al
JOIN 
    users u ON al.user_id = u.id
WHERE 
    al.credential_id = $1 AND al.group_id IS NULL
`

type GetAccessTypeAndUsersByCredentialIdRow struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	AccessType string    `json:"accessType"`
	PublicKey  string    `json:"publicKey"`
}

func (q *Queries) GetAccessTypeAndUsersByCredentialId(ctx context.Context, credentialID uuid.UUID) ([]GetAccessTypeAndUsersByCredentialIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getAccessTypeAndUsersByCredentialId, credentialID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAccessTypeAndUsersByCredentialIdRow{}
	for rows.Next() {
		var i GetAccessTypeAndUsersByCredentialIdRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.AccessType,
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

const getAllUrlsForUser = `-- name: GetAllUrlsForUser :many
SELECT DISTINCT
    field_value as value, credential_id as "credentialId"
FROM 
    fields
WHERE 
    user_id = $1 AND field_name = 'Domain'
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
    created_at,
    updated_at,
    name,
    description,
    folder_id,
    credential_type,
    created_by
FROM
    credentials
WHERE
    id = $1
`

type GetCredentialDataByIDRow struct {
	ID             uuid.UUID      `json:"id"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	Name           string         `json:"name"`
	Description    sql.NullString `json:"description"`
	FolderID       uuid.UUID      `json:"folderId"`
	CredentialType string         `json:"credentialType"`
	CreatedBy      uuid.UUID      `json:"createdBy"`
}

func (q *Queries) GetCredentialDataByID(ctx context.Context, id uuid.UUID) (GetCredentialDataByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getCredentialDataByID, id)
	var i GetCredentialDataByIDRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Description,
		&i.FolderID,
		&i.CredentialType,
		&i.CreatedBy,
	)
	return i, err
}

const getCredentialDetails = `-- name: GetCredentialDetails :one
SELECT
    id,
    NAME,
    COALESCE(description, '') AS "description"
FROM
    credentials
WHERE
    id = $1
`

type GetCredentialDetailsRow struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func (q *Queries) GetCredentialDetails(ctx context.Context, id uuid.UUID) (GetCredentialDetailsRow, error) {
	row := q.db.QueryRowContext(ctx, getCredentialDetails, id)
	var i GetCredentialDetailsRow
	err := row.Scan(&i.ID, &i.Name, &i.Description)
	return i, err
}

const getCredentialDetailsByIds = `-- name: GetCredentialDetailsByIds :many
SELECT
    C.id AS "credentialId",
    C.name,
    COALESCE(C.description, '') AS description,
    json_agg(
        json_build_object(
            'fieldName', COALESCE(ED.field_name, ''),
            'fieldValue', ED.field_value
        )
    ) AS "fields"
FROM
    credentials C
LEFT JOIN fields ED ON C.id = ED.credential_id AND ED.user_id = $2
WHERE
    C.id = ANY($1::UUID[])
GROUP BY C.id
`

type GetCredentialDetailsByIdsParams struct {
	Column1 []uuid.UUID `json:"column1"`
	UserID  uuid.UUID   `json:"userId"`
}

type GetCredentialDetailsByIdsRow struct {
	CredentialId uuid.UUID       `json:"credentialId"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Fields       json.RawMessage `json:"fields"`
}

func (q *Queries) GetCredentialDetailsByIds(ctx context.Context, arg GetCredentialDetailsByIdsParams) ([]GetCredentialDetailsByIdsRow, error) {
	rows, err := q.db.QueryContext(ctx, getCredentialDetailsByIds, pq.Array(arg.Column1), arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetCredentialDetailsByIdsRow{}
	for rows.Next() {
		var i GetCredentialDetailsByIdsRow
		if err := rows.Scan(
			&i.CredentialId,
			&i.Name,
			&i.Description,
			&i.Fields,
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
    c.id AS "credentialId"
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

const getCredentialsFieldsByIds = `-- name: GetCredentialsFieldsByIds :many
SELECT
    e.credential_id AS "credentialId",
    json_agg(
        json_build_object(
            'fieldId',
            e.id,
            'fieldValue',
            e.field_value
        )
    ) AS "fields"
FROM
    fields e
WHERE
    e.credential_id = ANY($1 :: uuid [ ])
    AND e.user_id = $2
GROUP BY
    e.credential_id
ORDER BY
    e.credential_id
`

type GetCredentialsFieldsByIdsParams struct {
	Column1 []uuid.UUID `json:"column1"`
	UserID  uuid.UUID   `json:"userId"`
}

type GetCredentialsFieldsByIdsRow struct {
	CredentialId uuid.UUID       `json:"credentialId"`
	Fields       json.RawMessage `json:"fields"`
}

func (q *Queries) GetCredentialsFieldsByIds(ctx context.Context, arg GetCredentialsFieldsByIdsParams) ([]GetCredentialsFieldsByIdsRow, error) {
	rows, err := q.db.QueryContext(ctx, getCredentialsFieldsByIds, pq.Array(arg.Column1), arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetCredentialsFieldsByIdsRow{}
	for rows.Next() {
		var i GetCredentialsFieldsByIdsRow
		if err := rows.Scan(&i.CredentialId, &i.Fields); err != nil {
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

const getEncryptedCredentialsByFolder = `-- name: GetEncryptedCredentialsByFolder :many
SELECT
    C .id as "credentialId",
    json_agg(
        json_build_object(
            'fieldName',
            e.field_name,
            'fieldValue',
            e.field_value
        )
    ) AS "encryptedFields"
FROM
    credentials C
    JOIN fields e ON C .id = e.credential_id
WHERE
    C .folder_id = $1
    AND e.user_id = $2
GROUP BY
    C .id
ORDER BY
    C .id
`

type GetEncryptedCredentialsByFolderParams struct {
	FolderID uuid.UUID `json:"folderId"`
	UserID   uuid.UUID `json:"userId"`
}

type GetEncryptedCredentialsByFolderRow struct {
	CredentialId    uuid.UUID       `json:"credentialId"`
	EncryptedFields json.RawMessage `json:"encryptedFields"`
}

func (q *Queries) GetEncryptedCredentialsByFolder(ctx context.Context, arg GetEncryptedCredentialsByFolderParams) ([]GetEncryptedCredentialsByFolderRow, error) {
	rows, err := q.db.QueryContext(ctx, getEncryptedCredentialsByFolder, arg.FolderID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetEncryptedCredentialsByFolderRow{}
	for rows.Next() {
		var i GetEncryptedCredentialsByFolderRow
		if err := rows.Scan(&i.CredentialId, &i.EncryptedFields); err != nil {
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
