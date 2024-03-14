-- sql/create_credential.sql
-- name: CreateCredential :one
INSERT INTO
    credentials (NAME, description, credential_type, folder_id, created_by, domain)
VALUES
    ($1, $2, $3, $4, $5, $6) RETURNING id;


-- name: GetCredentialDataByID :one
SELECT
    id,
    name,
    description,
    folder_id,
    credential_type,
    created_by,
    created_at,
    updated_at,
    updated_by
FROM
    credentials
WHERE
    id = $1;


-- name: FetchCredentialDetailsForUserByFolderId :many
SELECT
    C.id AS "credentialID",
    C.name,
    COALESCE(C.description, '') AS "description",
    C.credential_type AS "credentialType",
    C.created_at AS "createdAt",
    C.updated_at AS "updatedAt",
    C.created_by AS "createdBy",
    C.updated_by AS "updatedBy",
    A.access_type AS "accessType"
FROM
    credentials AS C,
    credential_access AS A
WHERE
    C.id = A .credential_id
    AND C.folder_id = $1
    AND A.user_id = $2;


-- name: GetAllUrlsForUser :many
SELECT DISTINCT
    field_value as value, credential_id as "credentialId"
FROM 
    fields
WHERE 
    user_id = $1 AND field_name = 'Domain';


-- name: GetCredentialDetailsByIDs :many
SELECT
    id,
    name,
    description,
    folder_id,
    credential_type,
    created_by,
    created_at,
    updated_at,
    updated_by
FROM
    credentials
WHERE
    id = ANY(@credentialIDs::UUID[]);



-- name: GetCredentialIdsByFolder :many
SELECT 
    c.id AS "credentialId"
FROM 
    credentials c
JOIN 
    credential_access a ON c.id = a.credential_id
WHERE 
    a.user_id = $1
    AND c.folder_id = $2;

-- name: GetAccessTypeAndUsersByCredentialId :many
SELECT 
    al.user_id as "id",
    u.name, 
    al.access_type,
    COALESCE(u.encryption_key, '') AS "publicKey"
FROM 
    credential_access al
JOIN 
    users u ON al.user_id = u.id
WHERE 
    al.credential_id = $1 AND al.group_id IS NULL;

-- name: GetAccessTypeAndGroupsByCredentialId :many
    SELECT DISTINCT
        al.group_id, 
        g.name,
        al.access_type
    FROM 
        credential_access al
    JOIN 
        groupings g ON al.group_id = g.id
    WHERE 
        al.credential_id = $1;


-- name: EditCredentialDetails :exec
UPDATE
    credentials
SET
    name = $2,
    description = $3,
    credential_type = $4,
    updated_at = NOW(),
    updated_by = $5
WHERE
    id = $1;

-- name: GetCredentialsForSearchByUserID :many
SELECT 
    c.id, 
    c.name, 
    COALESCE(c.description, '') AS description,
    COALESCE(c.domain, '') AS domain,
    c.folder_id, 
    COALESCE(f.name, '') AS folder_name
FROM 
    credentials c
JOIN 
    credential_access ca ON c.id = ca.credential_id
LEFT JOIN 
    folders f ON c.folder_id = f.id
WHERE 
    ca.user_id = $1;