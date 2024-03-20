-- name: AddFolderAccess :exec
INSERT INTO folder_access (folder_id, user_id, access_type, group_id)
VALUES ($1, $2, $3, $4);


-- name: CheckFolderAccessEntryExists :one
SELECT EXISTS (
    SELECT 1
    FROM folder_access
    WHERE user_id = $1 AND folder_id = $2 
    AND ((group_id IS NOT NULL AND group_id = $3) OR (group_id is null and $3 is null)) 
);


-- name: HasManageAccessForFolder :one
SELECT EXISTS (
  SELECT 1 FROM folder_access
  WHERE folder_id = $1 AND user_id = $2 AND access_type = 'manager'
);


-- name: HasReadAccessForFolder :one
SELECT EXISTS (
  SELECT 1 FROM folder_access
  WHERE folder_id = $1 AND user_id = $2 AND access_type IN ('manager', 'reader')
);


-- name: RemoveFolderAccessForUsers :exec
DELETE FROM folder_access WHERE group_id IS NULL 
AND folder_id = $1 AND user_id = ANY(@user_ids::UUID[]);


-- name: RemoveFolderAccessForGroups :exec
DELETE FROM folder_access WHERE folder_id = $1 AND group_id = ANY(@group_ids::UUID[]);


-- name: EditFolderAccessForUser :exec
UPDATE folder_access
SET access_type = $1
WHERE group_id IS NULL
AND folder_id = $2 AND user_id = $3;


-- name: EditFolderAccessForGroup :exec
UPDATE folder_access
SET access_type = $1
WHERE folder_id = $2 AND group_id = $3;


