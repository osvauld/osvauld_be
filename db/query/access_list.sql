-- name: AddToAccessList :one

INSERT INTO access_list (credential_id, user_id, access_type)
VALUES ($1, $2, $3)
RETURNING id;