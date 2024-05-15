
ALTER TABLE users ADD COLUMN created_by UUID;
-- SQL Definition for environments

CREATE TABLE environments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cli_user UUID NOT NULL REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updatedAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL REFERENCES users(id)
);
-- SQL Definition for environment fields

CREATE TABLE environment_fields (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    field_name VARCHAR(255) NOT NULL,
    field_value TEXT NOT NULL,
    parent_field_id UUID NOT NULL REFERENCES fields(id),
    env_id UUID NOT NULL REFERENCES environments(Id),
    cli_user UUID NOT NULL REFERENCES users(id),
    credential_id UUID NOT NULL REFERENCES credentials(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);