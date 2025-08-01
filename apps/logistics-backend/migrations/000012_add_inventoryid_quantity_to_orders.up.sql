-- Add inventory_id and quantity columns to orders
ALTER TABLE orders
ADD COLUMN inventory_id UUID NOT NULL,
ADD COLUMN quantity INTEGER NOT NULL CHECK (quantity > 0);

-- Add foreign key constraint linking orders.inventory_id -> inventories.id
ALTER TABLE orders
ADD CONSTRAINT fk_orders_inventory
FOREIGN KEY (inventory_id)
REFERENCES inventories(id)
ON DELETE CASCADE;