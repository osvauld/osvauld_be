-- name: CreateGroup :one
INSERT INTO groupings (name, created_by)
VALUES ($1, $2)
RETURNING id, name, created_by, created_at;

-- name: AddGroupMember :exec
INSERT INTO group_list (grouping_id, user_id, access_type)
VALUES ($1, $2, $3);

-- name: FetchUserGroups :many
SELECT groupings.id as "groupId", groupings.name, groupings.created_by, groupings.created_at
FROM groupings
JOIN group_list ON group_list.grouping_id = groupings.id
WHERE group_list.user_id = $1;

-- name: GetGroupMembers :many
SELECT users.id, users.name, users.username, COALESCE(users.encryption_key, '') as "publicKey"
FROM users
JOIN group_list ON users.id = group_list.user_id
WHERE group_list.grouping_id = $1;

-- name: CheckUserMemberOfGroup :one
SELECT EXISTS (
  SELECT 1 FROM group_list
  WHERE user_id = $1 AND grouping_id = $2
) as "exists";


-- name: CheckUserAdminOfGroup :one
SELECT EXISTS (
  SELECT 1 FROM group_list
  WHERE user_id = $1 AND grouping_id = $2 AND access_type = 'admin'
) as "exists";


-------------------------------------------------------------------------------------------------------


-- name: GetCredentialAccessDetailsWithGroupAccess :many
SELECT DISTINCT credential_id, access_type, folder_id
FROM credential_access
WHERE group_id = $1;


-- name: GetFolderIDAndTypeWithGroupAccess :many
SELECT DISTINCT folder_id, access_type
FROM folder_access
WHERE group_id = $1;

-- name: FetchCredentialIDsWithGroupAccess :many
SELECT distinct(credential_id) from credential_access
WHERE group_id = $1 and user_id = $2;


-- name: FetchCredentialAccessTypeForGroup :one
SELECT access_type FROM credential_access
WHERE group_id = $1 AND credential_id = $2;



-- name: FetchUsersByGroupIds :many
SELECT 
    g.id AS "groupId",
    json_agg(json_build_object('id', gl.user_id, 'publicKey', u.encryption_key)) AS "userDetails"
FROM 
    group_list gl
JOIN 
    groupings g ON gl.grouping_id = g.id
JOIN 
    users u ON gl.user_id = u.id
WHERE 
    g.id = ANY($1::UUID[])
GROUP BY 
    g.id;

-- name: GetGroupsWithoutAccess :many
SELECT id as "groupId", name 
FROM groupings
WHERE id NOT IN (
    SELECT group_id
    FROM folder_access
    WHERE folder_id = $1 AND group_id IS NOT NULL
);


-- name: GetUsersWithoutGroupAccess :many
SELECT u.id, u.username, u.name, COALESCE(u.encryption_key,'') as "encryptionKey"
FROM users u
LEFT JOIN group_list gl ON (u.id = gl.user_id AND gl.grouping_id = $1)
WHERE gl.grouping_id IS NULL  and u.status = 'active';



-- name: RemoveUserFromGroupList :exec
DELETE FROM group_list
WHERE user_id = $1 AND grouping_id = $2;


-- name: RemoveGroup :exec
DELETE FROM groupings
WHERE id = $1; 