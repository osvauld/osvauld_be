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
SELECT c.id AS credential_id, 
       c.name AS credential_name, 
       c.description AS credential_description, 
       json_agg(json_build_object('field_name', u.field_name, 'field_value', u.field_value)) AS unencrypted_data
FROM credentials c
JOIN access_list a ON c.id = a.credential_id
LEFT JOIN unencrypted_data u ON c.id = u.credential_id
WHERE a.user_id = $1 AND c.folder_id = $2
GROUP BY c.id, c.name, c.description;
----------------------------------------------------------------------

CREATE OR REPLACE FUNCTION share_secret(
    p_user_id UUID,
    p_credential_id UUID,
    p_field_names VARCHAR[],
    p_field_values VARCHAR[],
    p_access_type VARCHAR)
RETURNS VOID AS $$
DECLARE
    v_field_name VARCHAR;
    v_field_value VARCHAR;
BEGIN
    FOR i IN array_lower(p_field_names, 1)..array_upper(p_field_names, 1)
    LOOP
        v_field_name := p_field_names[i];
        v_field_value := p_field_values[i];

        INSERT INTO encrypted_data (user_id, credential_id, field_name, field_value)
        VALUES (p_user_id, p_credential_id, v_field_name, v_field_value);
    END LOOP;

    INSERT INTO access_list (user_id, credential_id, access_type)
    VALUES (p_user_id, p_credential_id, p_access_type);
END;
$$ LANGUAGE plpgsql;
----------------------------------------------------
-- name: ShareSecret :exec
SELECT share_secret($1, $2, $3, $4, $5);


-- name: GetCredentialDetails :one
SELECT id, name, description
FROM credentials
WHERE id = $1;


-- name: GetUserEncryptedData :many
SELECT field_name, field_value
FROM encrypted_data
WHERE user_id = $1 AND credential_id = $2;

-- name: GetCredentialUnencryptedData :many
SELECT field_name, field_value
FROM unencrypted_data
WHERE credential_id = $1;


-- name: AddCredential :one
SELECT add_credential_with_access($1::JSONB);
