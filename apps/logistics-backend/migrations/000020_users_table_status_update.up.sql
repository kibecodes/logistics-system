-- 0. Add the column allowing NULLs
ALTER TABLE users
ADD COLUMN status TEXT;

-- 1. Normalize existing data: replace NULLs and empty strings with 'active'
UPDATE users
SET status = 'active'
WHERE status IS NULL OR status = '';

-- 2. Set default for new rows
ALTER TABLE users
ALTER COLUMN status SET DEFAULT 'active';

-- 3. Disallow NULLs now that all existing rows are populated
ALTER TABLE users
ALTER COLUMN status SET NOT NULL;

-- 4. Add constraint for valid statuses
ALTER TABLE users
ADD CONSTRAINT users_status_check
CHECK (status IN ('active', 'inactive', 'suspended', 'pending'));
