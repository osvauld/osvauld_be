-- name: CreateGroup :one
INSERT INTO groups (name, members, created_by)
VALUES ($1, $2, $3)
RETURNING id;

-- name: AddMemberToGroup :exec
UPDATE groups 
SET members = array_cat(members, $3)
WHERE id = $1 AND created_by = $2
RETURNING id;
