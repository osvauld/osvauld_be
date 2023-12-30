-- name: FetchUnencryptedFieldsByCredentialID :many
SELECT
    id,
    field_name,
    field_value
FROM
    unencrypted_data
WHERE
    credential_id = $1;