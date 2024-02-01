-- name: AddCredentialAccess :one

INSERT INTO access_list (credential_id, user_id, access_type, group_id, folder_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;


-- name: GetCredentialAccessForUser :many
SELECT id, user_id, credential_id, group_id, access_type
FROM access_list
WHERE user_id = $1 AND credential_id = $2;


-- name: GetUsersByFolder :many
SELECT DISTINCT u.id, u.username, u.name, COALESCE(u.rsa_pub_key, '') as "publicKey"
FROM users u
JOIN access_list al ON u.id = al.user_id
JOIN credentials c ON al.credential_id = c.id
WHERE c.folder_id = $1;

-- name: GetUsersByCredential :many
SELECT users.id, users.username, users.name, COALESCE(users.rsa_pub_key, '') as "publicKey", access_list.access_type as "accessType"
FROM access_list
JOIN users ON access_list.user_id = users.id
WHERE access_list.credential_id = $1;


-- name: GetCredentialIDsByUserID :many
SELECT credential_id FROM access_list WHERE user_id = $1;


-- name: CheckCredentialAccessEntryExists :one
SELECT EXISTS (
    SELECT 1
    FROM access_list
    WHERE user_id = $1 AND credential_id = $2 
    AND ((group_id IS NOT NULL AND group_id = $3) OR (group_id is null and $3 is null)) 
    AND ((folder_id IS NOT NULL AND folder_id = $4) OR (folder_id is null and $4 is null))
);
