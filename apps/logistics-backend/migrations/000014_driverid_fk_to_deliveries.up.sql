-- Add foreign key constraint linking deliveries.driver_id -> users.id

ALTER TABLE deliveries
ADD CONSTRAINT drivers_user_fk
FOREIGN KEY (driver_id) 
REFERENCES users(id) 
ON DELETE CASCADE;