-- Step 1: Drop existing foreign key constraint
ALTER TABLE inventories DROP CONSTRAINT inventories_admin_id_fkey;

-- Step 2: Add it back with ON DELETE CASCADE
ALTER TABLE inventories
ADD CONSTRAINT fk_inventories_user
FOREIGN KEY (admin_id)
REFERENCES users(id)
ON DELETE CASCADE;