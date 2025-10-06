CREATE INDEX idx_orders_pickup_point
    ON orders USING GIST (pickup_point);

CREATE INDEX idx_orders_delivery_point
    ON orders USING GIST (delivery_point);
