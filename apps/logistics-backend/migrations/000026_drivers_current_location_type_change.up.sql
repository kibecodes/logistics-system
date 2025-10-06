-- 1. Drop and re-add the column as nullable geography
ALTER TABLE drivers
    DROP COLUMN IF EXISTS current_location,
    ADD COLUMN current_location GEOGRAPHY(Point, 4326);

-- 2. (Optional) Backfill existing rows with a default point, if needed
UPDATE drivers SET current_location = ST_GeogFromText('SRID=4326;POINT(0 0)') WHERE current_location IS NULL;

-- 3. Reapply the NOT NULL constraint once you're sure all rows have a value
ALTER TABLE drivers
    ALTER COLUMN current_location SET NOT NULL;

-- 4. Add the spatial index
CREATE INDEX IF NOT EXISTS idx_drivers_current_location
    ON drivers USING GIST (current_location);
