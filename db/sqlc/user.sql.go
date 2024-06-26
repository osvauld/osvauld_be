// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: user.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const checkCliUser = `-- name: CheckCliUser :one
SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND type = 'cli')
`

func (q *Queries) CheckCliUser(ctx context.Context, id uuid.UUID) (bool, error) {
	row := q.db.QueryRowContext(ctx, checkCliUser, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkIfUsersExist = `-- name: CheckIfUsersExist :one
SELECT EXISTS(SELECT 1 FROM users)
`

func (q *Queries) CheckIfUsersExist(ctx context.Context) (bool, error) {
	row := q.db.QueryRowContext(ctx, checkIfUsersExist)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkNameExist = `-- name: CheckNameExist :one
SELECT 
    EXISTS(SELECT 1 FROM users WHERE name = $1) AS exists
`

func (q *Queries) CheckNameExist(ctx context.Context, name string) (bool, error) {
	row := q.db.QueryRowContext(ctx, checkNameExist, name)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkUsernameExist = `-- name: CheckUsernameExist :one
SELECT 
    EXISTS(SELECT 1 FROM users WHERE username = $1) AS exists
`

func (q *Queries) CheckUsernameExist(ctx context.Context, username string) (bool, error) {
	row := q.db.QueryRowContext(ctx, checkUsernameExist, username)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
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

const createCliUser = `-- name: CreateCliUser :one
INSERT INTO users (

    username, 
    name, 
    encryption_key, 
    device_key, 
    temp_password, 
    registration_challenge, 
    signed_up, 
    type, 
    status, 
    created_by
) VALUES (
    $1, 
    $2, 
    $3, 
    $4, 
    $5, 
    $6, 
    $7, 
    $8, 
    $9, 
    $10 
)
RETURNING id
`

type CreateCliUserParams struct {
	Username              string         `json:"username"`
	Name                  string         `json:"name"`
	EncryptionKey         sql.NullString `json:"encryptionKey"`
	DeviceKey             sql.NullString `json:"deviceKey"`
	TempPassword          string         `json:"tempPassword"`
	RegistrationChallenge sql.NullString `json:"registrationChallenge"`
	SignedUp              bool           `json:"signedUp"`
	Type                  string         `json:"type"`
	Status                string         `json:"status"`
	CreatedBy             uuid.NullUUID  `json:"createdBy"`
}

func (q *Queries) CreateCliUser(ctx context.Context, arg CreateCliUserParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createCliUser,
		arg.Username,
		arg.Name,
		arg.EncryptionKey,
		arg.DeviceKey,
		arg.TempPassword,
		arg.RegistrationChallenge,
		arg.SignedUp,
		arg.Type,
		arg.Status,
		arg.CreatedBy,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (username, name, temp_password, type)
VALUES ($1, $2, $3, COALESCE($4, 'user'))
RETURNING id
`

type CreateUserParams struct {
	Username     string      `json:"username"`
	Name         string      `json:"name"`
	TempPassword string      `json:"tempPassword"`
	Column4      interface{} `json:"column4"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.Name,
		arg.TempPassword,
		arg.Column4,
	)
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

const getAllSignedUpUsers = `-- name: GetAllSignedUpUsers :many
SELECT id,name,username, COALESCE(encryption_key, '') AS "publicKey" FROM users where signed_up = true and type !='cli'
`

type GetAllSignedUpUsersRow struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	PublicKey string    `json:"publicKey"`
}

func (q *Queries) GetAllSignedUpUsers(ctx context.Context) ([]GetAllSignedUpUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllSignedUpUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllSignedUpUsersRow{}
	for rows.Next() {
		var i GetAllSignedUpUsersRow
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

const getAllUsers = `-- name: GetAllUsers :many
SELECT id,name,username, status, type FROM users
`

type GetAllUsersRow struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Username string    `json:"username"`
	Status   string    `json:"status"`
	Type     string    `json:"type"`
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
			&i.Status,
			&i.Type,
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

const getCliUsers = `-- name: GetCliUsers :many
SELECT id, username FROM users WHERE type = 'cli' and created_by = $1
`

type GetCliUsersRow struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

func (q *Queries) GetCliUsers(ctx context.Context, createdBy uuid.NullUUID) ([]GetCliUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, getCliUsers, createdBy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetCliUsersRow{}
	for rows.Next() {
		var i GetCliUsersRow
		if err := rows.Scan(&i.ID, &i.Username); err != nil {
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

const getRegistrationChallenge = `-- name: GetRegistrationChallenge :one
SELECT registration_challenge, status FROM users WHERE username = $1
`

type GetRegistrationChallengeRow struct {
	RegistrationChallenge sql.NullString `json:"registrationChallenge"`
	Status                string         `json:"status"`
}

func (q *Queries) GetRegistrationChallenge(ctx context.Context, username string) (GetRegistrationChallengeRow, error) {
	row := q.db.QueryRowContext(ctx, getRegistrationChallenge, username)
	var i GetRegistrationChallengeRow
	err := row.Scan(&i.RegistrationChallenge, &i.Status)
	return i, err
}

const getSuperUser = `-- name: GetSuperUser :one
select id, created_at, updated_at, username, name, encryption_key, device_key, temp_password, registration_challenge, signed_up, type, status, created_by from users where type = 'superadmin' limit 1
`

func (q *Queries) GetSuperUser(ctx context.Context) (User, error) {
	row := q.db.QueryRowContext(ctx, getSuperUser)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.Name,
		&i.EncryptionKey,
		&i.DeviceKey,
		&i.TempPassword,
		&i.RegistrationChallenge,
		&i.SignedUp,
		&i.Type,
		&i.Status,
		&i.CreatedBy,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, username, name, type FROM users WHERE id = $1
`

type GetUserByIDRow struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Name     string    `json:"name"`
	Type     string    `json:"type"`
}

func (q *Queries) GetUserByID(ctx context.Context, id uuid.UUID) (GetUserByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i GetUserByIDRow
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Name,
		&i.Type,
	)
	return i, err
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

const getUserDeviceKey = `-- name: GetUserDeviceKey :one
SELECT COALESCE(device_key,'') as "deviceKey" FROM users WHERE id = $1
`

func (q *Queries) GetUserDeviceKey(ctx context.Context, id uuid.UUID) (string, error) {
	row := q.db.QueryRowContext(ctx, getUserDeviceKey, id)
	var deviceKey string
	err := row.Scan(&deviceKey)
	return deviceKey, err
}

const getUserTempPassword = `-- name: GetUserTempPassword :one
SELECT temp_password, status FROM users WHERE username = $1
`

type GetUserTempPasswordRow struct {
	TempPassword string `json:"tempPassword"`
	Status       string `json:"status"`
}

func (q *Queries) GetUserTempPassword(ctx context.Context, username string) (GetUserTempPasswordRow, error) {
	row := q.db.QueryRowContext(ctx, getUserTempPassword, username)
	var i GetUserTempPasswordRow
	err := row.Scan(&i.TempPassword, &i.Status)
	return i, err
}

const getUserType = `-- name: GetUserType :one
SELECT type FROM users WHERE id = $1
`

func (q *Queries) GetUserType(ctx context.Context, id uuid.UUID) (string, error) {
	row := q.db.QueryRowContext(ctx, getUserType, id)
	var type_ string
	err := row.Scan(&type_)
	return type_, err
}

const removeUserFromOrg = `-- name: RemoveUserFromOrg :exec
DELETE FROM users WHERE id = $1
`

func (q *Queries) RemoveUserFromOrg(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, removeUserFromOrg, id)
	return err
}

const updateKeys = `-- name: UpdateKeys :exec
UPDATE users
SET encryption_key = $1, device_key = $2, signed_up = TRUE, status = 'active'
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

const updateRegistrationChallenge = `-- name: UpdateRegistrationChallenge :exec
UPDATE users
SET registration_challenge = $1, status = 'temp_login'
WHERE username = $2
`

type UpdateRegistrationChallengeParams struct {
	RegistrationChallenge sql.NullString `json:"registrationChallenge"`
	Username              string         `json:"username"`
}

func (q *Queries) UpdateRegistrationChallenge(ctx context.Context, arg UpdateRegistrationChallengeParams) error {
	_, err := q.db.ExecContext(ctx, updateRegistrationChallenge, arg.RegistrationChallenge, arg.Username)
	return err
}
