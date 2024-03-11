-- name: CreateUser :one
INSERT INTO users (username, name, temp_password)
VALUES ($1, $2, $3)
RETURNING id;


-- name: GetUserByUsername :one
SELECT id,name,username, COALESCE(encryption_key,'') as "publicKey"
FROM users
WHERE username = $1
LIMIT 1;


-- name: GetUserTempPassword :one
SELECT temp_password, status FROM users WHERE username = $1;

-- name: UpdateRegistrationChallenge :exec
UPDATE users
SET registration_challenge = $1, status = 'temp_login'
WHERE username = $2;


-- name: GetAllUsers :many
SELECT id,name,username, COALESCE(encryption_key, '') AS "publicKey" FROM users where signed_up = true;

-- name: GetUserByPublicKey :one
SELECT id
FROM users
WHERE device_key = $1
LIMIT 1;



-- name: CreateChallenge :one
INSERT INTO session_table (user_id, public_key, challenge)
VALUES ($1, $2, $3)
ON CONFLICT (public_key) DO UPDATE 
SET challenge = EXCLUDED.challenge,
    updated_at = CURRENT_TIMESTAMP
RETURNING *;


-- name: FetchChallenge :one
SELECT challenge FROM session_table WHERE user_id = $1;


-- name: UpdateKeys :one
UPDATE users
SET encryption_key = $1, device_key = $2, signed_up = TRUE, status = 'active'
WHERE username = $3
RETURNING id as "userId";



-- name: GetRegistrationChallenge :one
SELECT registration_challenge, status FROM users WHERE username = $1;