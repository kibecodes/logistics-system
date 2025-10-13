-- Drop redundant columns
ALTER TABLE inventories
    DROP COLUMN IF EXISTS location,
    DROP COLUMN IF EXISTS name,
    DROP COLUMN IF EXISTS slug,
    DROP COLUMN IF EXISTS admin_id;

-- Add new store_id column with NOT NULL and FK in one step
ALTER TABLE inventories
    ADD COLUMN store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE;

-- Optional: create index for store_id (for filtering/join performance)
CREATE INDEX IF NOT EXISTS idx_inventories_store_id ON inventories(store_id);
