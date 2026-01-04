-- Add join_code and is_locked to challenges table
ALTER TABLE challenges ADD COLUMN IF NOT EXISTS join_code VARCHAR(6) UNIQUE;
ALTER TABLE challenges ADD COLUMN IF NOT EXISTS is_locked BOOLEAN DEFAULT FALSE;

-- Backfill join codes for existing challenges
-- Generate unique 6-character alphanumeric codes
UPDATE challenges
SET join_code = UPPER(SUBSTRING(MD5(RANDOM()::TEXT || id::TEXT) FROM 1 FOR 6))
WHERE join_code IS NULL;

-- Make join_code NOT NULL after backfill
ALTER TABLE challenges ALTER COLUMN join_code SET NOT NULL;

-- Create index for fast lookup by join code
CREATE INDEX IF NOT EXISTS idx_challenges_join_code ON challenges(join_code);

-- Create index for locked challenges (for filtering)
CREATE INDEX IF NOT EXISTS idx_challenges_is_locked ON challenges(is_locked);
