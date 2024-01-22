-- sql/create_credential.sql
-- name: CreateCredential :one
INSERT INTO
    credentials (NAME, description, credential_type, folder_id, created_by)
VALUES
    ($1, $2, $3, $4, $5) RETURNING id;

-- name: CreateFieldData :one
INSERT INTO
    encrypted_data (field_name, field_value, credential_id, field_type, user_id)
VALUES
    ($1, $2, $3, $4, $5) RETURNING id;



-- name: FetchCredentialDataByID :one
SELECT
    id,
    created_at,
    updated_at,
    name,
    description,
    folder_id,
    created_by
FROM
    credentials
WHERE
    id = $1;

-- name: FetchCredentialsByUserAndFolder :many


WITH CredentialWithUnencrypted AS (
    SELECT
        C.id AS "id",
        C.name AS "name",
        COALESCE(C.description, '') AS "description",
        json_agg(
            json_build_object(
                'fieldName', u.field_name,
                'fieldValue', u.field_value
            )
        ) FILTER (WHERE u.field_name IS NOT NULL) AS "unencryptedFields"
    FROM
        credentials C
        LEFT JOIN unencrypted_data u ON C.id = u.credential_id
    WHERE
        C.folder_id = $2
    GROUP BY
        C.id
)
SELECT
    cwu.*
FROM
    CredentialWithUnencrypted cwu
JOIN
    access_list A ON cwu.id = A.credential_id
WHERE
    A.user_id = $1;

-- name: GetCredentialDetails :one
SELECT
    id,
    NAME,
    COALESCE(description, '') AS "description"
FROM
    credentials
WHERE
    id = $1;

-- name: GetUserEncryptedData :many
SELECT
    field_name AS "fieldName",
    field_value AS "fieldValue"
FROM
    encrypted_data
WHERE
    user_id = $1
    AND credential_id = $2;

-- name: GetCredentialUnencryptedData :many
SELECT
    field_name AS "fieldName",
    field_value AS "fieldValue"
FROM
    unencrypted_data
WHERE
    credential_id = $1;

-- name: AddCredential :one
SELECT
    add_credential_with_access($1 :: JSONB);

-- name: GetEncryptedCredentialsByFolder :many
SELECT
    C .id as "credentialId",
    json_agg(
        json_build_object(
            'fieldName',
            e.field_name,
            'fieldValue',
            e.field_value
        )
    ) AS "encryptedFields"
FROM
    credentials C
    JOIN encrypted_data e ON C .id = e.credential_id
WHERE
    C .folder_id = $1
    AND e.user_id = $2
GROUP BY
    C .id
ORDER BY
    C .id;

-- name: GetEncryptedDataByCredentialIds :many
SELECT
    e.credential_id AS "credentialId",
    json_agg(
        json_build_object(
            'fieldName',
            e.field_name,
            'fieldValue',
            e.field_value
        )
    ) AS "encryptedFields"
FROM
    encrypted_data e
WHERE
    e.credential_id = ANY($1 :: uuid [ ])
    AND e.user_id = $2
GROUP BY
    e.credential_id
ORDER BY
    e.credential_id;





-- -- name: GetCredentialsByUrl :many

-- WITH CredentialWithUnencrypted AS (
--     SELECT
--         C.id AS "id",
--         C.name AS "name",
--         COALESCE(C.description, '') AS "description",
--         json_agg(
--             json_build_object(
--                 'fieldName', u.field_name,
--                 'fieldValue', u.field_value,
--                 'isUrl', u.is_url,
--                 'url', u.url
--             )
--         ) FILTER (WHERE u.field_name IS NOT NULL) AS "unencryptedFields"
--     FROM
--         credentials C
--         LEFT JOIN unencrypted_data u ON C.id = u.credential_id
--     WHERE
--         C.id IN (SELECT credential_id FROM unencrypted_data as und WHERE und.url = $1)
--     GROUP BY
--         C.id
-- ),
-- DistinctAccess AS (
--     SELECT DISTINCT credential_id
--     FROM access_list
--     WHERE user_id = $2
-- )
-- SELECT
--     cwu.*
-- FROM
--     CredentialWithUnencrypted cwu
-- JOIN
--     DistinctAccess DA ON cwu.id = DA.credential_id;


-- name: GetAllUrlsForUser :many
SELECT DISTINCT
    COALESCE(ud.url, '') AS url
FROM 
    unencrypted_data ud
JOIN 
    credentials c ON ud.credential_id = c.id
JOIN 
    access_list al ON c.id = al.credential_id
WHERE 
    al.user_id = $1 AND ud.is_url = TRUE;

-- name: GetCredentialIdsByUrl :many
SELECT credential_id 
FROM unencrypted_data 
WHERE url = $1
AND credential_id IN (
    SELECT DISTINCT credential_id 
    FROM access_list 
    WHERE user_id = $2
);

-- name: GetCredentialDetailsByIds :many
SELECT
    C.id AS "credentialId",
    C.name,
    COALESCE(C.description, '') AS description,
    json_agg(
        json_build_object(
            'fieldName', COALESCE(ED.field_name, ''),
            'fieldValue', ED.field_value
        )
    ) AS "encryptedFields",
    json_agg(
        json_build_object(
            'fieldName', COALESCE(UD.field_name, ''),
            'fieldValue', UD.field_value,
            'isUrl', UD.is_url,
            'url', UD.url
        )
    ) AS "unencryptedFields"
FROM
    credentials C
LEFT JOIN encrypted_data ED ON C.id = ED.credential_id AND ED.user_id = $2
LEFT JOIN unencrypted_data UD ON C.id = UD.credential_id
WHERE
    C.id = ANY($1::UUID[])
GROUP BY C.id;
