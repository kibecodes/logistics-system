-- 1. Add currency column to payments
ALTER TABLE payments
ADD COLUMN currency CHAR(3) NOT NULL DEFAULT 'KES',
ADD CONSTRAINT currency_code_check CHECK (char_length(currency) = 3);

-- Change payments.amount to BIGINT cents
ALTER TABLE payments
ALTER COLUMN amount TYPE BIGINT
USING (amount * 100)::BIGINT;

-- 2. Inventories: replace price NUMERIC with amount + currency
ALTER TABLE inventories
RENAME COLUMN price TO old_price;

ALTER TABLE inventories
ADD COLUMN price_amount BIGINT NOT NULL DEFAULT 0,
ADD COLUMN price_currency CHAR(3) NOT NULL DEFAULT 'KES',
ADD CONSTRAINT price_currency_code_check CHECK (char_length(price_currency) = 3);

-- Migrate old decimal price to integer cents
UPDATE inventories
SET price_amount = (old_price * 100)::BIGINT,
    price_currency = 'KES';

-- Drop old column
ALTER TABLE inventories DROP COLUMN old_price;