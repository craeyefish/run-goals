-- Extend challenges table to support all goal types
-- Goal types: distance, elevation, summit_count, specific_summits

-- Add goal_type column
ALTER TABLE challenges ADD COLUMN IF NOT EXISTS goal_type VARCHAR(20) DEFAULT 'specific_summits';
ALTER TABLE challenges ADD CONSTRAINT check_goal_type
    CHECK (goal_type IN ('distance', 'elevation', 'summit_count', 'specific_summits'));

-- Add target_value for distance/elevation goals (stored in meters)
-- For distance: target in meters (will be displayed as km)
-- For elevation: target in meters
ALTER TABLE challenges ADD COLUMN IF NOT EXISTS target_value NUMERIC;

-- Rename target_count to target_summit_count for clarity
-- This is used for summit_count goal type
ALTER TABLE challenges RENAME COLUMN target_count TO target_summit_count;

-- Update existing challenges to have proper goal_type
-- Existing challenges are specific_summits type (already default)

-- For specific_summits: uses challenge_peaks table
-- For summit_count: uses target_summit_count
-- For distance/elevation: uses target_value
