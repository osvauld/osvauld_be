-- name: CreateCombinedField :exec
INSERT INTO combined_query_fields (user_id, combined_field)
VALUES ($1, '');