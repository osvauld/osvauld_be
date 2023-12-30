-- name: CreateUser :one
INSERT INTO users (username, name, public_key)
VALUES ($1, $2, $3)
RETURNING id;


-- name: GetUserByUsername :one
SELECT id,name,username, public_key as "publicKey"
FROM users
WHERE username = $1
LIMIT 1;


-- name: GetAllUsers :many
SELECT id,name,username, public_key AS "publicKey" FROM users;

-- name: GetUserByPublicKey :one
SELECT id
FROM users
WHERE ecc_pub_key = $1
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
