-- Pickup points
UPDATE orders
SET pickup_point = ST_SetSRID(ST_MakePoint(pickup_lng, pickup_lat), 4326)
WHERE pickup_point IS NULL;

-- Delivery points
UPDATE orders
SET delivery_point = ST_SetSRID(ST_MakePoint(delivery_lng, delivery_lat), 4326)
WHERE delivery_point IS NULL;
