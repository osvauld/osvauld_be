-- name: CreateCombinedField :exec
INSERT INTO combined_fields (user_id, combined_field)
VALUES ($1, '');