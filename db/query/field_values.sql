
-- name: AddFieldValue :exec
INSERT INTO field_values (field_id, field_value, user_id)
VALUES ($1, $2, $3);