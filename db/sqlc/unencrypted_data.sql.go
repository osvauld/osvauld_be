// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: unencrypted_data.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createUnencryptedData = `-- name: CreateUnencryptedData :one
INSERT INTO
    unencrypted_data (field_name, credential_id, field_value, is_url, url)
VALUES
    ($1, $2, $3, $4, $5) RETURNING id
`

type CreateUnencryptedDataParams struct {
	FieldName    string         `json:"fieldName"`
	CredentialID uuid.UUID      `json:"credentialId"`
	FieldValue   string         `json:"fieldValue"`
	IsUrl        bool           `json:"isUrl"`
	Url          sql.NullString `json:"url"`
}

func (q *Queries) CreateUnencryptedData(ctx context.Context, arg CreateUnencryptedDataParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createUnencryptedData,
		arg.FieldName,
		arg.CredentialID,
		arg.FieldValue,
		arg.IsUrl,
		arg.Url,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const fetchUnencryptedFieldsByCredentialID = `-- name: FetchUnencryptedFieldsByCredentialID :many
SELECT
    id,
    field_name,
    field_value
FROM
    unencrypted_data
WHERE
    credential_id = $1
`

type FetchUnencryptedFieldsByCredentialIDRow struct {
	ID         uuid.UUID `json:"id"`
	FieldName  string    `json:"fieldName"`
	FieldValue string    `json:"fieldValue"`
}

func (q *Queries) FetchUnencryptedFieldsByCredentialID(ctx context.Context, credentialID uuid.UUID) ([]FetchUnencryptedFieldsByCredentialIDRow, error) {
	rows, err := q.db.QueryContext(ctx, fetchUnencryptedFieldsByCredentialID, credentialID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []FetchUnencryptedFieldsByCredentialIDRow{}
	for rows.Next() {
		var i FetchUnencryptedFieldsByCredentialIDRow
		if err := rows.Scan(&i.ID, &i.FieldName, &i.FieldValue); err != nil {
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
