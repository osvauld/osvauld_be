// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: field_values.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const addFieldValue = `-- name: AddFieldValue :one
INSERT INTO field_values (field_id, field_value, user_id)
VALUES ($1, $2, $3) RETURNING id
`

type AddFieldValueParams struct {
	FieldID    uuid.UUID `json:"fieldId"`
	FieldValue string    `json:"fieldValue"`
	UserID     uuid.UUID `json:"userId"`
}

func (q *Queries) AddFieldValue(ctx context.Context, arg AddFieldValueParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, addFieldValue, arg.FieldID, arg.FieldValue, arg.UserID)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getFieldValueIDsForFieldIDs = `-- name: GetFieldValueIDsForFieldIDs :many
SELECT id, field_id FROM field_values WHERE field_id = ANY($2::UUID[]) and user_id = $1
`

type GetFieldValueIDsForFieldIDsParams struct {
	UserID   uuid.UUID   `json:"userId"`
	Fieldids []uuid.UUID `json:"fieldids"`
}

type GetFieldValueIDsForFieldIDsRow struct {
	ID      uuid.UUID `json:"id"`
	FieldID uuid.UUID `json:"fieldId"`
}

func (q *Queries) GetFieldValueIDsForFieldIDs(ctx context.Context, arg GetFieldValueIDsForFieldIDsParams) ([]GetFieldValueIDsForFieldIDsRow, error) {
	rows, err := q.db.QueryContext(ctx, getFieldValueIDsForFieldIDs, arg.UserID, pq.Array(arg.Fieldids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetFieldValueIDsForFieldIDsRow{}
	for rows.Next() {
		var i GetFieldValueIDsForFieldIDsRow
		if err := rows.Scan(&i.ID, &i.FieldID); err != nil {
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
