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