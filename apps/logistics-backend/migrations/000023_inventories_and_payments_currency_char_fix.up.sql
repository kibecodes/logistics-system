-- Migration: fix currency columns in inventories and payments

BEGIN;

-- ===== Inventories =====

-- 1. Ensure price_currency is VARCHAR(3) instead of CHAR(3)
ALTER TABLE inventories
ALTER COLUMN price_currency TYPE VARCHAR(3);

-- 2. Drop old constraints if they exist
ALTER TABLE inventories
DROP CONSTRAINT IF EXISTS price_currency_code_check;
ALTER TABLE inventories
DROP CONSTRAINT IF EXISTS price_currency_allowed_check;

-- 3. Add stricter regex check (must be exactly 3 uppercase letters)
ALTER TABLE inventories
ADD CONSTRAINT price_currency_code_check
CHECK (price_currency ~ '^[A-Z]{3}$');

-- 4. Optional: whitelist of supported currencies
ALTER TABLE inventories
ADD CONSTRAINT price_currency_allowed_check
CHECK (price_currency IN ('KES', 'USD', 'EUR', 'GBP'));


-- ===== Payments =====

-- 1. Ensure currency is VARCHAR(3) instead of CHAR(3)
ALTER TABLE payments
ALTER COLUMN currency TYPE VARCHAR(3);

-- 2. Drop old constraints if they exist
ALTER TABLE payments
DROP CONSTRAINT IF EXISTS currency_code_check;
ALTER TABLE payments
DROP CONSTRAINT IF EXISTS currency_allowed_check;

-- 3. Add stricter regex check
ALTER TABLE payments
ADD CONSTRAINT currency_code_check
CHECK (currency ~ '^[A-Z]{3}$');

-- 4. Optional: whitelist of supported currencies
ALTER TABLE payments
ADD CONSTRAINT currency_allowed_check
CHECK (currency IN ('KES', 'USD', 'EUR', 'GBP'));

COMMIT;
