
-- name: AddFieldValue :one
INSERT INTO field_values (field_id, field_value, user_id)
VALUES ($1, $2, $3) RETURNING id;

-- name: GetFieldValueIDsForFieldIDs :many
SELECT id, field_id FROM field_values WHERE field_id = ANY(@fieldIDs::UUID[]) and user_id = $1;