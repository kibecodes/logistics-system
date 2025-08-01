ALTER TABLE users 
DROP COLUMN IF EXISTS updated_at;

ALTER TABLE users 
DROP COLUMN IF EXISTS must_change_password;