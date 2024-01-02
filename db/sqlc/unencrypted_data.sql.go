// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: unencrypted_data.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

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
	FieldName  string    `json:"field_name"`
	FieldValue string    `json:"field_value"`
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
