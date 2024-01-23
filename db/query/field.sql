


-- name: FetchFieldNameAndTypeByFieldID :one
SELECT
    encrypted_data.field_name,
    encrypted_data.field_type
FROM encrypted_data
WHERE encrypted_data.id = $1;

