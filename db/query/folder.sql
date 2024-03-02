
-- name: AddFolder :one
INSERT INTO folders (name, description, created_by)
VALUES ($1, $2, $3)
RETURNING id, created_at;

-- name: AddFolderAccess :exec
INSERT INTO folder_access (folder_id, user_id, access_type, group_id)
VALUES ($1, $2, $3, $4);


-- name: FetchAccessibleFoldersForUser :many
SELECT id, name, description, created_at, created_by
FROM folders
WHERE id IN (
  SELECT DISTINCT(folder_id)
  FROM folder_access
  WHERE folder_access.user_id = $1
  UNION
  SELECT DISTINCT(c.folder_id)
  FROM credentials as c
  JOIN credential_access as a ON c.id = a.credential_id
  WHERE a.user_id = $1
);

-- name: CheckFolderAccessEntryExists :one
SELECT EXISTS (
    SELECT 1
    FROM folder_access
    WHERE user_id = $1 AND folder_id = $2 
    AND ((group_id IS NOT NULL AND group_id = $3) OR (group_id is null and $3 is null)) 
);




-- name: IsFolderOwner :one
SELECT EXISTS (
  SELECT 1 FROM folder_access
  WHERE folder_id = $1 AND user_id = $2 AND access_type = 'owner'
);


-- name: GetSharedUsersForFolder :many
SELECT users.id, users.name, users.username, COALESCE(users.encryption_key,'') as "publicKey", folder_access.access_type as "accessType"
FROM folder_access
JOIN users ON folder_access.user_id = users.id
WHERE folder_access.folder_id = $1;

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
