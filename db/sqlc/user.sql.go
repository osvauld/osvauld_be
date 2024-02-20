// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: user.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const checkTempPassword = `-- name: CheckTempPassword :one
SELECT COUNT(*) FROM users WHERE username = $1 AND temp_password = $2
`

type CheckTempPasswordParams struct {
	Username     string `json:"username"`
	TempPassword string `json:"tempPassword"`
}

func (q *Queries) CheckTempPassword(ctx context.Context, arg CheckTempPasswordParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, checkTempPassword, arg.Username, arg.TempPassword)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createChallenge = `-- name: CreateChallenge :one
INSERT INTO session_table (user_id, public_key, challenge)
VALUES ($1, $2, $3)
ON CONFLICT (public_key) DO UPDATE 
SET challenge = EXCLUDED.challenge,
    updated_at = CURRENT_TIMESTAMP
RETURNING id, user_id, public_key, challenge, device_id, session_id, created_at, updated_at
`

type CreateChallengeParams struct {
	UserID    uuid.UUID `json:"userId"`
	PublicKey string    `json:"publicKey"`
	Challenge string    `json:"challenge"`
}

func (q *Queries) CreateChallenge(ctx context.Context, arg CreateChallengeParams) (SessionTable, error) {
	row := q.db.QueryRowContext(ctx, createChallenge, arg.UserID, arg.PublicKey, arg.Challenge)
	var i SessionTable
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.PublicKey,
		&i.Challenge,
		&i.DeviceID,
		&i.SessionID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (username, name, temp_password)
VALUES ($1, $2, $3)
RETURNING id
`

type CreateUserParams struct {
	Username     string `json:"username"`
	Name         string `json:"name"`
	TempPassword string `json:"tempPassword"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.Name, arg.TempPassword)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const fetchChallenge = `-- name: FetchChallenge :one
SELECT challenge FROM session_table WHERE user_id = $1
`

func (q *Queries) FetchChallenge(ctx context.Context, userID uuid.UUID) (string, error) {
	row := q.db.QueryRowContext(ctx, fetchChallenge, userID)
	var challenge string
	err := row.Scan(&challenge)
	return challenge, err
}

const getAllUsers = `-- name: GetAllUsers :many
SELECT id,name,username, COALESCE(encryption_key, '') AS "publicKey" FROM users where signed_up = true
`

type GetAllUsersRow struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	PublicKey string    `json:"publicKey"`
}

func (q *Queries) GetAllUsers(ctx context.Context) ([]GetAllUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllUsersRow{}
	for rows.Next() {
		var i GetAllUsersRow
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

const getUserByPublicKey = `-- name: GetUserByPublicKey :one
SELECT id
FROM users
WHERE device_key = $1
LIMIT 1
`

func (q *Queries) GetUserByPublicKey(ctx context.Context, deviceKey sql.NullString) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getUserByPublicKey, deviceKey)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id,name,username, COALESCE(encryption_key,'') as "publicKey"
FROM users
WHERE username = $1
LIMIT 1
`

type GetUserByUsernameRow struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	PublicKey string    `json:"publicKey"`
}

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (GetUserByUsernameRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i GetUserByUsernameRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Username,
		&i.PublicKey,
	)
	return i, err
}

const updateKeys = `-- name: UpdateKeys :exec
UPDATE users
SET encryption_key = $1, device_key = $2, signed_up = TRUE
WHERE username = $3
`

type UpdateKeysParams struct {
	EncryptionKey sql.NullString `json:"encryptionKey"`
	DeviceKey     sql.NullString `json:"deviceKey"`
	Username      string         `json:"username"`
}

func (q *Queries) UpdateKeys(ctx context.Context, arg UpdateKeysParams) error {
	_, err := q.db.ExecContext(ctx, updateKeys, arg.EncryptionKey, arg.DeviceKey, arg.Username)
	return err
}
