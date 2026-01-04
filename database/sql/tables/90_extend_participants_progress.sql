-- Extend challenge_participants table for tracking progress across all goal types

-- Add progress fields for distance goals (stored in meters)
ALTER TABLE challenge_participants ADD COLUMN IF NOT EXISTS total_distance NUMERIC DEFAULT 0;

-- Add progress fields for elevation goals (stored in meters)
ALTER TABLE challenge_participants ADD COLUMN IF NOT EXISTS total_elevation NUMERIC DEFAULT 0;

-- Add progress fields for summit_count goals
ALTER TABLE challenge_participants ADD COLUMN IF NOT EXISTS total_summit_count INTEGER DEFAULT 0;

-- Keep existing fields:
-- - peaks_completed: for specific_summits goals
-- - total_peaks: for specific_summits goals

-- Create indexes for performance on progress queries
CREATE INDEX IF NOT EXISTS idx_participants_distance ON challenge_participants(total_distance);
CREATE INDEX IF NOT EXISTS idx_participants_elevation ON challenge_participants(total_elevation);
CREATE INDEX IF NOT EXISTS idx_participants_summit_count ON challenge_participants(total_summit_count);

-- These fields will be updated by:
-- 1. ProcessActivityForChallenges service when activities are synced
-- 2. RefreshParticipantProgress when manually triggered
