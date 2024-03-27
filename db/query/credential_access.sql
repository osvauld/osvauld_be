-- name: AddCredentialAccess :one
INSERT INTO credential_access (credential_id, user_id, access_type, group_id, folder_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;


-- name: GetCredentialAccessTypeForUser :many
SELECT id, user_id, credential_id, group_id, access_type
FROM credential_access
WHERE user_id = $1 AND credential_id = $2;


-- name: HasManageAccessForCredential :one
SELECT EXISTS (
  SELECT 1 FROM credential_access
  WHERE credential_id = $1 AND user_id = $2 AND access_type = 'manager'
);


-- name: HasReadAccessForCredential :one
SELECT EXISTS (
  SELECT 1 FROM credential_access
  WHERE credential_id = $1 AND user_id = $2 AND access_type IN ('reader', 'manager')
);


-- name: GetUsersByCredential :many
SELECT users.id, users.username, users.name, COALESCE(users.encryption_key, '') as "publicKey", credential_access.access_type as "accessType"
FROM credential_access
JOIN users ON credential_access.user_id = users.id
WHERE credential_access.credential_id = $1;


-- name: GetCredentialIDsByUserID :many
SELECT credential_id FROM credential_access WHERE user_id = $1;


-- name: CheckCredentialAccessEntryExists :one
SELECT EXISTS (
    SELECT 1
    FROM credential_access
    WHERE user_id = $1 AND credential_id = $2 
    AND ((group_id IS NOT NULL AND group_id = $3) OR (group_id is null and $3 is null)) 
    AND ((folder_id IS NOT NULL AND folder_id = $4) OR (folder_id is null and $4 is null))
);


-- name: RemoveCredentialAccessForUsers :exec
DELETE FROM credential_access WHERE group_id IS NULL AND folder_id IS NULL 
AND credential_id = $1 AND user_id = ANY(@user_ids::UUID[]);


-- name: RemoveCredentialAccessForUsersWithFolderID :exec
DELETE FROM credential_access WHERE group_id IS NULL
AND folder_id = $1 AND user_id = ANY(@user_ids::UUID[]);


-- name: RemoveCredentialAccessForGroups :exec
DELETE FROM credential_access WHERE folder_id IS NULL 
AND credential_id = $1 AND group_id = ANY(@group_ids::UUID[]);


-- name: RemoveCredentialAccessForGroupsWithFolderID :exec
DELETE FROM credential_access WHERE folder_id = $1 AND group_id = ANY(@group_ids::UUID[]);


-- name: EditCredentialAccessForUser :exec
UPDATE credential_access
SET access_type = $1
WHERE  group_id IS NULL AND folder_id IS NULL
AND credential_id = $2 AND user_id = $3;


-- name: EditCredentialAccessForGroupWithFolderID :exec
UPDATE credential_access
SET access_type = $1
WHERE folder_id = $2 AND group_id = $3;


-- name: EditCredentialAccessForUserWithFolderID :exec
UPDATE credential_access
SET access_type = $1
WHERE group_id IS NULL
AND folder_id = $2 AND user_id = $3;


-- name: EditCredentialAccessForGroup :exec
UPDATE credential_access
SET access_type = $1
WHERE folder_id IS NULL
AND credential_id = $2 AND group_id = $3;


-- name: GetCredentialUsersWithDirectAccess :many
SELECT 
    ca.user_id,
    u.name,
    u.username,
    ca.access_type,
    CASE WHEN ca.folder_id IS NULL THEN 'acquired' ELSE 'inherited' END AS "accessSource"
FROM 
    credential_access ca
JOIN 
    users u ON ca.user_id = u.id
WHERE 
    ca.credential_id = $1 AND ca.group_id IS NULL;


-- name: GetCredentialGroups :many
SELECT 
    ca.group_id,
    g.name,
    ca.access_type,
    CASE WHEN ca.folder_id IS NULL THEN 'acquired' ELSE 'inherited' END AS "accessSource"
FROM 
    credential_access ca
JOIN 
    groupings g ON g.id = ca.id
WHERE 
    ca.credential_id = $1;


-- name: GetCredentialUsersForDataSync :many
SELECT 
    DISTINCT ca.user_id as "id",
    COALESCE(u.encryption_key, '') AS "publicKey"
FROM
    credential_access ca
JOIN
    users u ON ca.user_id = u.id
WHERE
    ca.credential_id = $1;

