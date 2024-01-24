


-- name: FetchFieldNameAndTypeByFieldIDForUser :one
SELECT
    encrypted_data.field_name,
    encrypted_data.field_type
FROM encrypted_data
WHERE encrypted_data.id = $1 AND encrypted_data.user_id = $2;

