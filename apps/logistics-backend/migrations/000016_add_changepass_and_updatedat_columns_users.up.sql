-- Add Updated_at column to users table
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS must_change_password BOOLEAN DEFAULT FALSE;

ALTER TABLE users 
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT now();