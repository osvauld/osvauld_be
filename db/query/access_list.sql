-- name: AddToAccessList :one

INSERT INTO access_list (credential_id, user_id, access_type)
VALUES ($1, $2, $3)
RETURNING id;


-- name: HasUserAccess :one
SELECT EXISTS (
  SELECT 1
  FROM access_list
  WHERE user_id = $1 AND credential_id = $2
) AS has_access;