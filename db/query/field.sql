

-- name: AddFieldData :one
INSERT INTO
    encrypted_data (field_name, field_value, credential_id, field_type, user_id)
VALUES
    ($1, $2, $3, $4, $5) RETURNING id;


-- name: GetFieldDataByCredentialIDsForUser :many
SELECT
    encrypted_data.id as "fieldId",
    encrypted_data.credential_id as "credentialId",
    encrypted_data.field_name,
    encrypted_data.field_value,
    encrypted_data.field_type
FROM encrypted_data
WHERE encrypted_data.user_id = $1 
AND encrypted_data.credential_id = ANY(@Credentials::UUID[]);