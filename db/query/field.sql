


-- name: FetchFieldNameAndTypeByFieldIDForUser :one
SELECT
    encrypted_data.field_name,
    encrypted_data.field_type
FROM encrypted_data
WHERE encrypted_data.id = $1 AND encrypted_data.user_id = $2;


-- name: AddFieldData :one
INSERT INTO
    encrypted_data (field_name, field_value, credential_id, field_type, user_id)
VALUES
    ($1, $2, $3, $4, $5) RETURNING id;
