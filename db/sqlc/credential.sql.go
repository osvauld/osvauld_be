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

const addCredential = `-- name: AddCredential :one
SELECT
    add_credential_with_access($1 :: JSONB)
`

func (q *Queries) AddCredential(ctx context.Context, dollar_1 json.RawMessage) (interface{}, error) {
	row := q.db.QueryRowContext(ctx, addCredential, dollar_1)
	var add_credential_with_access interface{}
	err := row.Scan(&add_credential_with_access)
	return add_credential_with_access, err
}

const createCredential = `-- name: CreateCredential :one
INSERT INTO
    credentials (NAME, description, credential_type, folder_id, created_by)
VALUES
    ($1, $2, $3, $4, $5) RETURNING id
`

type CreateCredentialParams struct {
	Name           string         `json:"name"`
	Description    sql.NullString `json:"description"`
	CredentialType string         `json:"credential_type"`
	FolderID       uuid.UUID      `json:"folder_id"`
	CreatedBy      uuid.UUID      `json:"created_by"`
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

const createFieldData = `-- name: CreateFieldData :one
INSERT INTO
    encrypted_data (field_name, field_value, credential_id, field_type, user_id)
VALUES
    ($1, $2, $3, $4, $5) RETURNING id
`

type CreateFieldDataParams struct {
	FieldName    string    `json:"field_name"`
	FieldValue   string    `json:"field_value"`
	CredentialID uuid.UUID `json:"credential_id"`
	FieldType    string    `json:"field_type"`
	UserID       uuid.UUID `json:"user_id"`
}

func (q *Queries) CreateFieldData(ctx context.Context, arg CreateFieldDataParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createFieldData,
		arg.FieldName,
		arg.FieldValue,
		arg.CredentialID,
		arg.FieldType,
		arg.UserID,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const fetchCredentialDataByID = `-- name: FetchCredentialDataByID :one
SELECT
    id,
    created_at,
    updated_at,
    name,
    description,
    folder_id,
    created_by
FROM
    credentials
WHERE
    id = $1
`

type FetchCredentialDataByIDRow struct {
	ID          uuid.UUID      `json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	FolderID    uuid.UUID      `json:"folder_id"`
	CreatedBy   uuid.UUID      `json:"created_by"`
}

func (q *Queries) FetchCredentialDataByID(ctx context.Context, id uuid.UUID) (FetchCredentialDataByIDRow, error) {
	row := q.db.QueryRowContext(ctx, fetchCredentialDataByID, id)
	var i FetchCredentialDataByIDRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Description,
		&i.FolderID,
		&i.CreatedBy,
	)
	return i, err
}

const fetchCredentialsByUserAndFolder = `-- name: FetchCredentialsByUserAndFolder :many


WITH CredentialWithUnencrypted AS (
    SELECT
        C.id AS "id",
        C.name AS "name",
        COALESCE(C.description, '') AS "description",
        json_agg(
            json_build_object(
                'fieldName', u.field_name,
                'fieldValue', u.field_value
            )
        ) FILTER (WHERE u.field_name IS NOT NULL) AS "unencryptedFields"
    FROM
        credentials C
        LEFT JOIN unencrypted_data u ON C.id = u.credential_id
    WHERE
        C.folder_id = $2
    GROUP BY
        C.id
)
SELECT
    cwu.id, cwu.name, cwu.description, cwu."unencryptedFields"
FROM
    CredentialWithUnencrypted cwu
JOIN
    access_list A ON cwu.id = A.credential_id
WHERE
    A.user_id = $1
`

type FetchCredentialsByUserAndFolderParams struct {
	UserID   uuid.UUID `json:"user_id"`
	FolderID uuid.UUID `json:"folder_id"`
}

type FetchCredentialsByUserAndFolderRow struct {
	ID                uuid.UUID       `json:"id"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	UnencryptedFields json.RawMessage `json:"unencryptedFields"`
}

func (q *Queries) FetchCredentialsByUserAndFolder(ctx context.Context, arg FetchCredentialsByUserAndFolderParams) ([]FetchCredentialsByUserAndFolderRow, error) {
	rows, err := q.db.QueryContext(ctx, fetchCredentialsByUserAndFolder, arg.UserID, arg.FolderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []FetchCredentialsByUserAndFolderRow{}
	for rows.Next() {
		var i FetchCredentialsByUserAndFolderRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.UnencryptedFields,
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
    COALESCE(ud.url, '') AS url
FROM 
    unencrypted_data ud
JOIN 
    credentials c ON ud.credential_id = c.id
JOIN 
    access_list al ON c.id = al.credential_id
WHERE 
    al.user_id = $1 AND ud.is_url = TRUE
`

// -- name: GetCredentialsByUrl :many
// WITH CredentialWithUnencrypted AS (
//
//	SELECT
//	    C.id AS "id",
//	    C.name AS "name",
//	    COALESCE(C.description, '') AS "description",
//	    json_agg(
//	        json_build_object(
//	            'fieldName', u.field_name,
//	            'fieldValue', u.field_value,
//	            'isUrl', u.is_url,
//	            'url', u.url
//	        )
//	    ) FILTER (WHERE u.field_name IS NOT NULL) AS "unencryptedFields"
//	FROM
//	    credentials C
//	    LEFT JOIN unencrypted_data u ON C.id = u.credential_id
//	WHERE
//	    C.id IN (SELECT credential_id FROM unencrypted_data as und WHERE und.url = $1)
//	GROUP BY
//	    C.id
//
// ),
// DistinctAccess AS (
//
//	SELECT DISTINCT credential_id
//	FROM access_list
//	WHERE user_id = $2
//
// )
// SELECT
//
//	cwu.*
//
// FROM
//
//	CredentialWithUnencrypted cwu
//
// JOIN
//
//	DistinctAccess DA ON cwu.id = DA.credential_id;
func (q *Queries) GetAllUrlsForUser(ctx context.Context, userID uuid.UUID) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getAllUrlsForUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		items = append(items, url)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
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
    ) AS "encryptedFields",
    json_agg(
        json_build_object(
            'fieldName', COALESCE(UD.field_name, ''),
            'fieldValue', UD.field_value,
            'isUrl', UD.is_url,
            'url', UD.url
        )
    ) AS "unencryptedFields"
FROM
    credentials C
LEFT JOIN encrypted_data ED ON C.id = ED.credential_id AND ED.user_id = $2
LEFT JOIN unencrypted_data UD ON C.id = UD.credential_id
WHERE
    C.id = ANY($1::UUID[])
GROUP BY C.id
`

type GetCredentialDetailsByIdsParams struct {
	Column1 []uuid.UUID `json:"column_1"`
	UserID  uuid.UUID   `json:"user_id"`
}

type GetCredentialDetailsByIdsRow struct {
	CredentialId      uuid.UUID       `json:"credentialId"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	EncryptedFields   json.RawMessage `json:"encryptedFields"`
	UnencryptedFields json.RawMessage `json:"unencryptedFields"`
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
			&i.EncryptedFields,
			&i.UnencryptedFields,
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

const getCredentialIdsByUrl = `-- name: GetCredentialIdsByUrl :many
SELECT credential_id 
FROM unencrypted_data 
WHERE url = $1
AND credential_id IN (
    SELECT DISTINCT credential_id 
    FROM access_list 
    WHERE user_id = $2
)
`

type GetCredentialIdsByUrlParams struct {
	Url    sql.NullString `json:"url"`
	UserID uuid.UUID      `json:"user_id"`
}

func (q *Queries) GetCredentialIdsByUrl(ctx context.Context, arg GetCredentialIdsByUrlParams) ([]uuid.UUID, error) {
	rows, err := q.db.QueryContext(ctx, getCredentialIdsByUrl, arg.Url, arg.UserID)
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

const getCredentialUnencryptedData = `-- name: GetCredentialUnencryptedData :many
SELECT
    field_name AS "fieldName",
    field_value AS "fieldValue"
FROM
    unencrypted_data
WHERE
    credential_id = $1
`

type GetCredentialUnencryptedDataRow struct {
	FieldName  string `json:"fieldName"`
	FieldValue string `json:"fieldValue"`
}

func (q *Queries) GetCredentialUnencryptedData(ctx context.Context, credentialID uuid.UUID) ([]GetCredentialUnencryptedDataRow, error) {
	rows, err := q.db.QueryContext(ctx, getCredentialUnencryptedData, credentialID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetCredentialUnencryptedDataRow{}
	for rows.Next() {
		var i GetCredentialUnencryptedDataRow
		if err := rows.Scan(&i.FieldName, &i.FieldValue); err != nil {
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
    JOIN encrypted_data e ON C .id = e.credential_id
WHERE
    C .folder_id = $1
    AND e.user_id = $2
GROUP BY
    C .id
ORDER BY
    C .id
`

type GetEncryptedCredentialsByFolderParams struct {
	FolderID uuid.UUID `json:"folder_id"`
	UserID   uuid.UUID `json:"user_id"`
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

const getEncryptedDataByCredentialIds = `-- name: GetEncryptedDataByCredentialIds :many
SELECT
    e.credential_id AS "credentialId",
    json_agg(
        json_build_object(
            'fieldName',
            e.field_name,
            'fieldValue',
            e.field_value
        )
    ) AS "encryptedFields"
FROM
    encrypted_data e
WHERE
    e.credential_id = ANY($1 :: uuid [ ])
    AND e.user_id = $2
GROUP BY
    e.credential_id
ORDER BY
    e.credential_id
`

type GetEncryptedDataByCredentialIdsParams struct {
	Column1 []uuid.UUID `json:"column_1"`
	UserID  uuid.UUID   `json:"user_id"`
}

type GetEncryptedDataByCredentialIdsRow struct {
	CredentialId    uuid.UUID       `json:"credentialId"`
	EncryptedFields json.RawMessage `json:"encryptedFields"`
}

func (q *Queries) GetEncryptedDataByCredentialIds(ctx context.Context, arg GetEncryptedDataByCredentialIdsParams) ([]GetEncryptedDataByCredentialIdsRow, error) {
	rows, err := q.db.QueryContext(ctx, getEncryptedDataByCredentialIds, pq.Array(arg.Column1), arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetEncryptedDataByCredentialIdsRow{}
	for rows.Next() {
		var i GetEncryptedDataByCredentialIdsRow
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

const getUserEncryptedData = `-- name: GetUserEncryptedData :many
SELECT
    field_name AS "fieldName",
    field_value AS "fieldValue"
FROM
    encrypted_data
WHERE
    user_id = $1
    AND credential_id = $2
`

type GetUserEncryptedDataParams struct {
	UserID       uuid.UUID `json:"user_id"`
	CredentialID uuid.UUID `json:"credential_id"`
}

type GetUserEncryptedDataRow struct {
	FieldName  string `json:"fieldName"`
	FieldValue string `json:"fieldValue"`
}

func (q *Queries) GetUserEncryptedData(ctx context.Context, arg GetUserEncryptedDataParams) ([]GetUserEncryptedDataRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserEncryptedData, arg.UserID, arg.CredentialID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetUserEncryptedDataRow{}
	for rows.Next() {
		var i GetUserEncryptedDataRow
		if err := rows.Scan(&i.FieldName, &i.FieldValue); err != nil {
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
