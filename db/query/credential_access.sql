-- name: AddCredentialAccess :one

INSERT INTO credential_access (credential_id, user_id, access_type, group_id, folder_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;


-- name: GetCredentialAccessForUser :many
SELECT id, user_id, credential_id, group_id, access_type
FROM credential_access
WHERE user_id = $1 AND credential_id = $2;


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