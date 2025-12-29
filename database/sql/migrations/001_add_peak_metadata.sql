-- Migration: Add metadata columns to peaks table for better differentiation
-- Run this manually against your database before deploying the updated backend

-- Add new columns
ALTER TABLE peaks ADD COLUMN IF NOT EXISTS alt_name VARCHAR;
ALTER TABLE peaks ADD COLUMN IF NOT EXISTS name_en VARCHAR;
ALTER TABLE peaks ADD COLUMN IF NOT EXISTS region VARCHAR;
ALTER TABLE peaks ADD COLUMN IF NOT EXISTS wikipedia VARCHAR;
ALTER TABLE peaks ADD COLUMN IF NOT EXISTS wikidata VARCHAR;
ALTER TABLE peaks ADD COLUMN IF NOT EXISTS description TEXT;
ALTER TABLE peaks ADD COLUMN IF NOT EXISTS prominence NUMERIC;

-- Optional: Create index on region for filtering
CREATE INDEX IF NOT EXISTS idx_peaks_region ON peaks(region);

-- Optional: Create index on name for search
CREATE INDEX IF NOT EXISTS idx_peaks_name ON peaks(name);
