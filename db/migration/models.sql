-- SQL Definition for BaseModel (common fields)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


-- SQL Definition for User
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    username VARCHAR(255) UNIQUE NOT NULL
);
-- SQL Definition for Folder
CREATE TABLE folders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    created_by UUID REFERENCES users(id)
);

-- SQL Definition for Credential
CREATE TABLE credentials (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    folder_id UUID REFERENCES folders(id),
    created_by UUID REFERENCES users(id)
);


-- SQL Definition for AccessList
CREATE TABLE access_list (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    credential_id UUID REFERENCES credentials(id),
    user_id UUID REFERENCES users(id),
    access_type VARCHAR(255) NOT NULL
);


-- SQL Definition for EncryptedData
CREATE TABLE encrypted_data (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    field_name VARCHAR(255) NOT NULL,
    credential_id UUID REFERENCES credentials(id),
    field_value VARCHAR(255) NOT NULL,
    user_id UUID REFERENCES users(id)
);

-- SQL Definition for UnencryptedData
CREATE TABLE unencrypted_data (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    field_name VARCHAR(255) NOT NULL,
    credential_id UUID REFERENCES credentials(id),
    field_value VARCHAR(255) NOT NULL
);



-- SQL Definition for Group
CREATE TABLE groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    name VARCHAR(255) NOT NULL,
    members UUID[] NOT NULL,
    created_by UUID REFERENCES users(id)
);
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
CREATE OR REPLACE FUNCTION add_credential_with_access(
    jsonb_input JSONB
) RETURNS UUID AS $$
DECLARE
    v_credential_id UUID;
    v_folder_id UUID;
    v_name TEXT;
    v_description TEXT;
    v_unencrypted_fields JSONB;
    v_encrypted_fields JSONB;
    v_unique_user_ids JSONB;
    v_user_id UUID;
    v_unencrypted_field JSONB;
    v_encrypted_field JSONB;
    v_encrypted_field_data JSONB;
    v_created_by UUID;
BEGIN
    -- Extract fields from input
    v_name := jsonb_input->>'name';
    v_description := jsonb_input->>'description';
    v_folder_id := (jsonb_input->>'folderId')::UUID;
    v_unencrypted_fields := jsonb_input->'unencryptedFields';
    v_encrypted_fields := jsonb_input->'encryptedFields';
    v_unique_user_ids := jsonb_input->'uniqueUserIds';
    v_created_by := (jsonb_input->>'createdBy')::UUID;

    -- Create the credential
    INSERT INTO credentials (name, description, folder_id, created_by)
    VALUES (v_name, v_description, v_folder_id, v_created_by)
    RETURNING id INTO v_credential_id;

    -- Add unencrypted fields
    FOR v_unencrypted_field IN SELECT * FROM jsonb_array_elements(v_unencrypted_fields)
    LOOP
        INSERT INTO unencrypted_data (field_name, field_value, credential_id)
        VALUES ((v_unencrypted_field->>'fieldName')::varchar(255), (v_unencrypted_field->>'fieldValue')::varchar(255), v_credential_id);
    END LOOP;

    -- Add encrypted fields and access list entries
    FOR v_encrypted_field_data IN SELECT * FROM jsonb_array_elements(v_encrypted_fields)
    LOOP
        v_user_id := (v_encrypted_field_data->>'userId')::UUID;
        FOR v_encrypted_field IN SELECT * FROM jsonb_array_elements(v_encrypted_field_data->'fields')
        LOOP
            INSERT INTO encrypted_data (field_name, field_value, credential_id, user_id)
            VALUES ((v_encrypted_field->>'fieldName')::varchar(255), (v_encrypted_field->>'fieldValue')::varchar(255), v_credential_id, v_user_id);
        END LOOP;
    END LOOP;

    -- Process unique user IDs for access list
    FOR v_user_id IN SELECT * FROM jsonb_array_elements_text(v_unique_user_ids)
    LOOP
        INSERT INTO access_list (credential_id, user_id, access_type)
        VALUES (v_credential_id, v_user_id::UUID, 'default_access');
    END LOOP;

    RETURN v_credential_id;
END;
$$ LANGUAGE plpgsql;