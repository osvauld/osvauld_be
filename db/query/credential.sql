-- sql/create_credential.sql
-- name: CreateCredential :one
INSERT INTO credentials (name, description, folder_id, created_by)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: CreateEncryptedData :one
INSERT INTO encrypted_data (field_name, credential_id, field_value, user_id)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: CreateUnencryptedData :one
INSERT INTO unencrypted_data (field_name, credential_id, field_value)
VALUES ($1, $2, $3)
RETURNING id;

-- name: FetchCredentialsByUserAndFolder :many
SELECT
  c.id AS "id",
  c.name AS "name",
  COALESCE(c.description, '') AS "description",
  json_agg(
    json_build_object(
      'fieldName', u.field_name,  -- Assuming actual column name is in snake_case
      'fieldValue', u.field_value  -- Assuming actual column name is in snake_case
    )
  ) AS "unencryptedData"
FROM credentials c
JOIN access_list a ON c.id = a.credential_id
LEFT JOIN unencrypted_data u ON c.id = u.credential_id
WHERE a.user_id = $1 AND c.folder_id = $2
GROUP BY c.id;
-- name: ShareSecret :exec
SELECT share_secret($1, $2, $3, $4, $5);


-- name: GetCredentialDetails :one
SELECT id, name, COALESCE(description, '') AS "description" 
FROM credentials
WHERE id = $1;


-- name: GetUserEncryptedData :many
SELECT field_name AS fieldName, field_value AS fieldValue
FROM encrypted_data
WHERE user_id = $1 AND credential_id = $2;

-- name: GetCredentialUnencryptedData :many
SELECT field_name AS fieldName, field_value AS fieldValue
FROM unencrypted_data
WHERE credential_id = $1;


-- name: AddCredential :one
SELECT add_credential_with_access($1::JSONB);
