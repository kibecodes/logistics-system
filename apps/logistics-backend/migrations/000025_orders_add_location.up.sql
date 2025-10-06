-- 002_orders_add_location.sql
ALTER TABLE orders
    ADD COLUMN pickup_lat DOUBLE PRECISION,
    ADD COLUMN pickup_lng DOUBLE PRECISION,
    ADD COLUMN delivery_lat DOUBLE PRECISION,
    ADD COLUMN delivery_lng DOUBLE PRECISION;

ALTER TABLE orders
    ADD COLUMN pickup_point GEOGRAPHY(Point, 4326),
    ADD COLUMN delivery_point GEOGRAPHY(Point, 4326);
