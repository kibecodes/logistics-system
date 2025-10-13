-- Rollback: drop the foreign key and column, restore old structure placeholders
ALTER TABLE inventories
    DROP CONSTRAINT IF EXISTS fk_inventories_store,
    DROP COLUMN IF EXISTS store_id,
    ADD COLUMN admin_id UUID,
    ADD COLUMN name TEXT,
    ADD COLUMN slug TEXT,
    ADD COLUMN location TEXT;

DROP INDEX IF EXISTS idx_inventories_store_id;