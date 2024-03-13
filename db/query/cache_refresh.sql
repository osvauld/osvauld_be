
-- name: CreateCacheRefresh :exec
INSERT INTO cache_refresh (folder_id, user_id, credential_id, type)
VALUES ($1, $2, $3, $4);


-- name: GetCredentialIdsByUserIdForCacheRefresh :many
SELECT credential_id AS "credentialId"
FROM cache_refresh
WHERE user_id = $1;