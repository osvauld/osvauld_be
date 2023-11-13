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