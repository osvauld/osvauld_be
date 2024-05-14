-- name: AddEnvironment :one
INSERT INTO environments (
    cli_user, 
    name, 
    created_by
) VALUES (
    $1, 
    $2, 
    $3
)
RETURNING Id;


-- name: CheckCredentialExistsForEnv :one
SELECT EXISTS (
    SELECT 1 
    FROM environment_fields 
    WHERE credential_id = $1 AND env_id = $2
);


-- name: CreateEnvFields :one
INSERT INTO environment_fields (
    cli_user, 
    credential_id, 
    field_value, 
    field_name, 
    parent_field_id, 
    env_id
) VALUES (
    $1, 
    $2, 
    $3, 
    $4, 
    $5, 
    $6
)
RETURNING id;

-- name: GetEnvironmentsForUser :many
SELECT e.*,   COALESCE( u.encryption_key, '') as "publicKey"
FROM environments e
JOIN users u ON e.cli_user = u.id
WHERE e.cli_user IN (
    SELECT id 
    FROM users 
    WHERE u.created_by = $1 AND type = 'cli'
);


-- name: GetEnvironmentByID :one
SELECT * from environments WHERE id = $1 and created_by = $2;

-- name: GetEnvFields :many
SELECT pf.field_value, ef.field_name, ef.id ,ef.credential_id 
FROM environment_fields ef
JOIN fields f ON ef.parent_field_id = f.id
JOIN fields pf ON ef.parent_field_id = pf.id
WHERE ef.env_id = $1;


