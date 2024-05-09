
ALTER TABLE users ADD COLUMN created_by UUID;
-- SQL Definition for environments

CREATE TABLE environments (
    Id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cli_user UUID REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    createdAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updatedAt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID
);
-- SQL Definition for environment fields

CREATE TABLE environment_fields (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    field_name VARCHAR(255) NOT NULL,
    field_value TEXT,
    parent_field_id UUID REFERENCES fields(id),
    field_id UUID REFERENCES fields(id),
    env_id UUID REFERENCES environments(Id),
    credential_id UUID REFERENCES credentials(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);