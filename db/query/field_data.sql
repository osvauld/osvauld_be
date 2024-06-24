-- name: AddFieldData :one
INSERT INTO field_data (field_name, field_type, credential_id, created_by, updated_by)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: DeleteFields :exec 
DELETE FROM field_data WHERE id = ANY(@field_ids::UUID[]);