
-- name: AddFolder :one
INSERT INTO folders (name, description, created_by)
VALUES ($1, $2, $3)
RETURNING id, created_at;


-- name: FetchAccessibleFoldersForUser :many
SELECT folders.id, folders.name, folders.description, folders.created_at, folders.created_by, COALESCE(folder_access.access_type, 'none') as "accessType"
FROM folders
LEFT JOIN folder_access ON folders.id = folder_access.folder_id AND folder_access.user_id = $1
WHERE folders.id IN (
  SELECT DISTINCT(folder_id)
  FROM folder_access
  WHERE folder_access.user_id = $1
  UNION
  SELECT DISTINCT(c.folder_id)
  FROM credentials as c
  JOIN credential_access as a ON c.id = a.credential_id
  WHERE a.user_id = $1
);

-- name: GetSharedGroupsForFolder :many
SELECT g.id, g.name, f.access_type
FROM folder_access AS f JOIN groupings AS g ON f.group_id = g.id
WHERE f.folder_id = $1
group by g.id, g.name, f.access_type;


-- name: GetFolderAccessForUser :many
SELECT access_type FROM folder_access
WHERE folder_id = $1 AND user_id = $2;


-- name: GetAccessTypeAndUserByFolder :many
SELECT user_id, access_type
FROM folder_access
WHERE folder_id = $1;

-- name: IsUserManagerOrOwner :one
SELECT EXISTS (
  SELECT 1 FROM folder_access
  WHERE folder_id = $1 AND user_id = $2 AND access_type IN ('owner', 'manager')
);

-- name: RemoveFolder :exec
DELETE FROM folders
WHERE id = $1;

-- name: EditFolder :exec
UPDATE folders
SET name = $2, description = $3
WHERE id = $1;