-- sql/create_credential.sql
-- name: CreateCredential :one
INSERT INTO
    credentials (NAME, description, folder_id, created_by)
VALUES
    ($1, $2, $3, $4) RETURNING id;

-- name: CreateEncryptedData :one
INSERT INTO
    encrypted_data (field_name, credential_id, field_value, user_id)
VALUES
    ($1, $2, $3, $4) RETURNING id;

-- name: CreateUnencryptedData :one
INSERT INTO
    unencrypted_data (field_name, credential_id, field_value)
VALUES
    ($1, $2, $3) RETURNING id;


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
