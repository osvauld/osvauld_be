-- SQL Definition for BaseModel (common fields)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


-- SQL Definition for User
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    username VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL, 
    public_key TEXT NOT NULL 
);
-- SQL Definition for Folder
CREATE TABLE folders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    created_by UUID NOT NULL REFERENCES users(id)
);

-- SQL Definition for Credential
CREATE TABLE credentials (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    folder_id UUID NOT NULL REFERENCES folders(id),
    created_by UUID NOT NULL REFERENCES users(id)
);


-- SQL Definition for AccessList
CREATE TABLE access_list (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    credential_id UUID NOT NULL REFERENCES credentials(id),
    user_id UUID NOT NULL REFERENCES users(id),
    access_type VARCHAR(255) NOT NULL,
    group_id UUID REFERENCES groupings(id)
);


-- SQL Definition for EncryptedData
CREATE TABLE encrypted_data (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    field_name VARCHAR(255) NOT NULL,
    credential_id UUID NOT NULL REFERENCES credentials(id),
    field_value TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id)
);

-- SQL Definition for UnencryptedData
CREATE TABLE unencrypted_data (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    field_name VARCHAR(255) NOT NULL,
    credential_id UUID NOT NULL REFERENCES credentials(id),
    field_value VARCHAR(255) NOT NULL
);



-- SQL Definition for Group
CREATE TABLE groupings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    name VARCHAR(255) NOT NULL,
    created_by UUID REFERENCES users(id)
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


-- SQL Definition for Folder Access
CREATE TABLE folder_access (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    folder_id UUID NOT NULL REFERENCES folders(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    access_type VARCHAR(255) NOT NULL,
    UNIQUE(folder_id, user_id)
);


CREATE OR REPLACE FUNCTION share_secret(
    jsonb_input JSONB
) RETURNS VOID AS $$
DECLARE
    v_user_id UUID;
    v_credential_id UUID;
    v_field_names TEXT[];
    v_field_values TEXT[];
    v_access_type VARCHAR;
    v_field_name VARCHAR;
    v_field_value TEXT;
BEGIN
    -- Extract fields from input
    v_user_id := (jsonb_input->>'userId')::UUID;
    v_credential_id := (jsonb_input->>'credentialId')::UUID;
    v_field_names := ARRAY(SELECT jsonb_array_elements_text(jsonb_input->'fieldNames'));
    v_field_values := ARRAY(SELECT jsonb_array_elements_text(jsonb_input->'fieldValues'));
    v_access_type := jsonb_input->>'accessType';

    FOR i IN array_lower(v_field_names, 1)..array_upper(v_field_names, 1)
    LOOP
        v_field_name := v_field_names[i];
        v_field_value := v_field_values[i];

        INSERT INTO encrypted_data (user_id, credential_id, field_name, field_value)
        VALUES (v_user_id, v_credential_id, v_field_name, v_field_value);
    END LOOP;

    INSERT INTO access_list (user_id, credential_id, access_type)
    VALUES (v_user_id, v_credential_id, v_access_type);
END;
$$ LANGUAGE plpgsql;
