-- name: FetchUnencryptedFieldsByCredentialID :many
SELECT
    id,
    field_name,
    field_value
FROM
    unencrypted_data
WHERE
    credential_id = $1;



-- name: CreateUnencryptedData :one
INSERT INTO
    unencrypted_data (field_name, credential_id, field_value, is_url, url)
VALUES
    ($1, $2, $3, $4, $5) RETURNING id;