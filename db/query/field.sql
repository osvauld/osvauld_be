

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


-- name: RemoveCredentialFieldsForUsers :exec
DELETE FROM fields WHERE credential_id = $1 AND user_id = ANY(@user_ids::UUID[]);


-- name: DeleteAccessRemovedFields :exec
DELETE FROM fields
WHERE
    EXISTS (
        -- Select fields rows that don't have a corresponding entry in credential_access
        SELECT 1
        FROM fields f
        WHERE
            NOT EXISTS (
                -- Look for a matching entry in credential_access
                SELECT 1
                FROM credential_access ca
                WHERE
                    ca.credential_id = f.credential_id
                    AND ca.user_id = f.user_id
            )
            AND f.credential_id = fields.credential_id
            AND f.user_id = fields.user_id
    );



