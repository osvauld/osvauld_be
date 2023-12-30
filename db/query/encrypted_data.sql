-- name: FetchEncryptedFieldsByCredentialIDAndUserID :many
SELECT
    id,
    field_name,
    field_value
FROM
    encrypted_data
WHERE
    credential_id = $1
    AND user_id = $2;