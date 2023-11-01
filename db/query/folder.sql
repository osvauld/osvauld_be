-- name: CreateFolder :one
INSERT INTO folders (name, description, created_by)
VALUES ($1, $2, $3)
RETURNING id;


-- name: FetchAccessibleAndCreatedFoldersByUser :many
WITH unique_credential_ids AS (
  SELECT DISTINCT credential_id
  FROM access_list
  WHERE user_id = $1
),
unique_folder_ids AS (
  SELECT DISTINCT folder_id
  FROM credentials
  WHERE id IN (SELECT credential_id FROM unique_credential_ids)
)
SELECT f.*
FROM folders f
WHERE f.id IN (SELECT folder_id FROM unique_folder_ids)
   OR f.created_by = $1;
