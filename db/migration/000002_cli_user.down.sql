
-- Down command for users
ALTER TABLE users DROP COLUMN IF EXISTS created_by CASCADE;

-- Down command for environments
DROP TABLE IF EXISTS environments CASCADE;

-- Down command for environment_fields
DROP TABLE IF EXISTS environment_fields CASCADE;