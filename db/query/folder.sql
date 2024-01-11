-- name: CreateFolder :one
WITH new_folder AS (
  INSERT INTO folders (name, description, created_by)
  VALUES ($1, $2, $3)
  RETURNING id
),
folder_access_insert AS (
  INSERT INTO folder_access (folder_id, user_id, access_type)
  SELECT id, $3, 'owner' FROM new_folder
)
SELECT id FROM new_folder;

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
SELECT 
    id, 
    name, 
    COALESCE(description, '') AS description 
FROM folders f
WHERE f.id IN (SELECT folder_id FROM unique_folder_ids)
   OR f.created_by = $1;

-- name: IsFolderOwner :one
SELECT EXISTS (
  SELECT 1 FROM folder_access
  WHERE folder_id = $1 AND user_id = $2 AND access_type = 'owner'
);
-- name: AddFolderAccess :exec
INSERT INTO folder_access (folder_id, user_id, access_type)
SELECT $1, unnest($2::uuid[]), unnest($3::text[]);


-- name: GetSharedUsers :many
SELECT users.id, users.name, users.username, COALESCE(users.rsa_pub_key,'') as "publicKey", folder_access.access_type as "accessType"
FROM folder_access
JOIN users ON folder_access.user_id = users.id
WHERE folder_access.folder_id = $1;


-- name: GetFolderAccessForUser :many
SELECT access_type FROM folder_access
WHERE folder_id = $1 AND user_id = $2;



-- name: GetAccessTypeAndUserByFolder :many
SELECT user_id, access_type
FROM folder_access
WHERE folder_id = $1;