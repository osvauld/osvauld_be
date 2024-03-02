

-- name: AddField :one
INSERT INTO
    fields (field_name, field_value, credential_id, field_type, user_id, created_by)
VALUES
    ($1, $2, $3, $4, $5, $6) RETURNING id;


-- name: DeleteCredentialFields :exec
DELETE FROM fields
WHERE credential_id = $1;

-- name: GetNonSensitiveFieldsForCredentialIDs :many
SELECT
    f.id,
    f.credential_id,
    f.field_name,
    f.field_value,
    f.field_type
FROM fields as f
WHERE
field_type != 'sensitive' 
AND f.user_id = $1 
AND f.credential_id = ANY(@credentials::UUID[]);


-- name: GetAllFieldsForCredentialIDs :many
SELECT
    f.id,
    f.credential_id,
    f.field_name,
    f.field_value,
    f.field_type
FROM fields as f
WHERE
f.user_id = $1 
AND f.credential_id = ANY(@credentials::UUID[]);


-- name: CheckFieldEntryExists :one
SELECT EXISTS (
    SELECT 1
    FROM fields
    WHERE credential_id = $1 AND user_id = $2
);

-- name: GetSensitiveFields :many
SELECT
    f.id,
    f.field_name,
    f.field_value
FROM fields as f
WHERE
field_type = 'sensitive'
AND f.credential_id = $1
AND f.user_id = $2;
