-- SQL Definition for BaseModel (common fields)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


-- SQL Definition for User
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    username VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL, 
    encryption_key TEXT,
    device_key TEXT,
    temp_password VARCHAR(255) NOT NULL,
    registration_challenge VARCHAR(255),
    signed_up BOOLEAN NOT NULL DEFAULT FALSE,
    type VARCHAR(255) NOT NULL DEFAULT 'user',
    status VARCHAR(255) NOT NULL DEFAULT 'created'
);
-- SQL Definition for Folder
CREATE TABLE folders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    name VARCHAR(255) NOT NULL,
    description VARCHAR(2048),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL
);

-- SQL Definition for Credential
CREATE TABLE credentials (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description VARCHAR(2048),
    credential_type VARCHAR(255) NOT NULL,
    folder_id UUID NOT NULL REFERENCES folders(id) ON DELETE CASCADE,
    domain VARCHAR(2048),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- SQL Definition for Fields
CREATE TABLE fields (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    field_name VARCHAR(255) NOT NULL,
    field_value TEXT NOT NULL,
    field_type VARCHAR(255) NOT NULL,
    credential_id UUID NOT NULL REFERENCES credentials(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID  REFERENCES users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL
);



-- SQL Definition for Group
CREATE TABLE groupings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    name VARCHAR(255) NOT NULL,
    created_by UUID  REFERENCES users(id) ON DELETE SET NULL
);


-- SQL Definition for Group List
CREATE TABLE group_list (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    grouping_id UUID NOT NULL REFERENCES groupings(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    access_type VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(grouping_id, user_id)
);


-- SQL Definition for AccessList
CREATE TABLE credential_access (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    credential_id UUID NOT NULL REFERENCES credentials(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    access_type VARCHAR(255) NOT NULL,
    group_id UUID REFERENCES groupings(id) ON DELETE CASCADE,
    folder_id UUID REFERENCES folders(id) ON DELETE CASCADE
);


-- SQL Definition for Folder Access
CREATE TABLE folder_access (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    folder_id UUID NOT NULL REFERENCES folders(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    access_type VARCHAR(255) NOT NULL,
    group_id UUID REFERENCES groupings(id) ON DELETE CASCADE
);


-- SQL Definition for session table

CREATE TABLE session_table (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    public_key TEXT NOT NULL UNIQUE,
    challenge VARCHAR(255) NOT NULL,
    device_id VARCHAR(255),
    session_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_user_id ON session_table(user_id);
CREATE INDEX idx_session_id ON session_table(session_id);

