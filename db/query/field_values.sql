
-- name: AddFieldValue :one
INSERT INTO field_values (field_id, field_value, user_id)
VALUES ($1, $2, $3) RETURNING id;