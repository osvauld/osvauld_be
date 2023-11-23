-- name: CreateGroup :exec
WITH new_group AS (
  INSERT INTO groupings (name, created_by)
  VALUES ($1, $2)
  RETURNING id
)
INSERT INTO group_list (grouping_id, user_id, access_type)
SELECT id, $2, 'owner'
FROM new_group;

-- name: AddMemberToGroup :exec

INSERT INTO group_list (grouping_id, user_id, access_type)
SELECT $1, u.id, 'member'
FROM unnest($2::uuid[]) AS u(id)
LEFT JOIN group_list gl ON gl.grouping_id = $1 AND gl.user_id = u.id
WHERE gl.user_id IS NULL;
-- name: GetUserGroups :many
SELECT g.*
FROM groupings g
JOIN group_list gl ON g.id = gl.grouping_id
WHERE gl.user_id = $1;



-- name: GetGroupMembers :many
SELECT u.id, u.name, u.username, u.public_key as "publicKey"
FROM users u
JOIN group_list gl ON u.id = gl.user_id
WHERE gl.grouping_id = $1;
