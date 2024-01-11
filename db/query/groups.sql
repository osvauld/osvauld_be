

-- name: CreateGroup :one
INSERT INTO groupings (name, created_by)
VALUES ($1, $2)
RETURNING id;


-- name: AddGroupMemberRecord :exec
INSERT INTO group_list (grouping_id, user_id, access_type)
VALUES ($1, $2, $3);


-- name: GetGroupMembers :many
SELECT users.id, users.name, users.username, COALESCE(users.rsa_pub_key, '') as "publicKey"
FROM users
JOIN group_list ON users.id = group_list.user_id
WHERE group_list.grouping_id = $1;

-- name: FetchUserGroups :many
SELECT groupings.id, groupings.name, groupings.created_by, groupings.created_at
FROM groupings
JOIN group_list ON group_list.grouping_id = groupings.id
WHERE group_list.user_id = $1;

-- name: CheckUserMemberOfGroup :one
SELECT EXISTS (
  SELECT 1 FROM group_list
  WHERE user_id = $1 AND grouping_id = $2
) as "exists";


-- name: FetchGroupAccessType :one
SELECT access_type FROM group_list
WHERE user_id = $1 AND grouping_id = $2;

-- name: FetchCredentialIDsWithGroupAccess :many
SELECT distinct(credential_id) from access_list
WHERE group_id = $1;


-- name: FetchCredentialAccessTypeForGroupMember :one
SELECT access_type FROM access_list
WHERE group_id = $1 AND credential_id = $2 AND user_id = $3;