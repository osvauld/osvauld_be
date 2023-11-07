-- name: CreateUser :one
INSERT INTO users (username, name, public_key)
VALUES ($1, $2, $3)
RETURNING id;


-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1
LIMIT 1;