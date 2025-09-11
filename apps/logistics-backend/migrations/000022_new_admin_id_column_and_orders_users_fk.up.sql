-- 1. Add the new column, nullable for now
ALTER TABLE orders
ADD COLUMN admin_id UUID;

-- 2. Add the foreign key constraint
ALTER TABLE orders
ADD CONSTRAINT fk_orders_admin_id_users_id
FOREIGN KEY (admin_id) REFERENCES users (id);

-- 3. Backfill old rows with a valid admin_id 
UPDATE orders SET admin_id = '7435c977-65e0-4e8f-9634-26a66059f47c' WHERE admin_id IS NULL;

-- 4. Enforce NOT NULL going forward (do this AFTER backfilling, otherwise it will fail)
ALTER TABLE orders
ALTER COLUMN admin_id SET NOT NULL;