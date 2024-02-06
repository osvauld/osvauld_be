

-- name: AddFieldData :one
INSERT INTO
    fields (field_name, field_value, credential_id, field_type, user_id)
VALUES
    ($1, $2, $3, $4, $5) RETURNING id;


-- name: GetFieldDataByCredentialIDsForUser :many
SELECT
    f.id,
    f.credential_id,
    f.field_name,
    f.field_value,
    f.field_type
FROM fields as f
WHERE f.user_id = $1 
AND f.credential_id = ANY(@credentials::UUID[]);


-- name: CheckFieldEntryExists :one
SELECT EXISTS (
    SELECT 1
    FROM fields
    WHERE credential_id = $1 AND user_id = $2
);

-- name: FetchEncryptedFieldsByCredentialIDAndUserID :many
SELECT
    id,
    field_name,
    field_value
FROM
    fields
WHERE
    credential_id = $1
    AND user_id = $2;
