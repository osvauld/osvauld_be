-- name: GetCredentialIDsByUserID :many
SELECT credential_id FROM access_list WHERE user_id = $1;
