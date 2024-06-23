

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
    fd.id,
    fd.credential_id,
    fd.field_name,
    fv.field_value,
    fd.field_type
FROM field_data as fd
JOIN field_values as fv ON fd.id = fv.field_id
WHERE
(field_type = 'meta' OR field_type = 'additional')
AND fv.user_id = $1 
AND fd.field_type != 'totp'
AND fd.credential_id = ANY(@credentials::UUID[]);


-- name: GetAllFieldsForCredentialIDs :many
SELECT
    fd.id,
    fd.field_name,
    fv.field_value,
    fd.field_type,
    fd.credential_id
FROM field_data as fd
JOIN field_values as fv ON fd.id = fv.field_id
WHERE
fv.user_id = $1 
AND fd.credential_id = ANY(@credentials::UUID[]);


-- name: CheckFieldEntryExists :one
SELECT EXISTS (
    SELECT 1
    FROM fields
    WHERE credential_id = $1 AND user_id = $2
);

-- name: GetSensitiveFields :many
SELECT
    fd.id,
    fd.field_name,
    fv.field_value,
    fd.field_type
FROM field_data as fd
JOIN field_values as fv ON fd.id = fv.field_id
WHERE
(fd.field_type = 'sensitive' OR fd.field_type = 'totp')
AND fd.credential_id = $1
AND fv.user_id = $2;


-- name: RemoveCredentialFieldsForUsers :exec
DELETE FROM fields WHERE credential_id = $1 AND user_id = ANY(@user_ids::UUID[]);


-- name: DeleteAccessRemovedFields :exec
DELETE FROM field_values fv
WHERE EXISTS (
    -- Select fields_values rows that don't have a corresponding entry in credential_access
    SELECT 1
    FROM field_data fd
    WHERE
        fd.id = fv.field_id
        AND NOT EXISTS (
            -- Look for a matching entry in credential_access
            SELECT 1
            FROM credential_access ca
            WHERE
                ca.credential_id = fd.credential_id
                AND ca.user_id = fv.user_id
        )
);


-- name: EditFieldData :exec
UPDATE field_data
SET field_name = $1, field_type = $2, updated_by = $3, updated_at = NOW()
WHERE id = $4;


-- name: EditFieldValue :exec
UPDATE field_values
SET field_value = $1
WHERE field_id = $2 AND user_id = $3;

