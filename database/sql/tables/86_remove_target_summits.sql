-- Remove target_summits column from personal_yearly_goals
-- This data is now stored in the summit_favourites table (year-independent)
ALTER TABLE personal_yearly_goals DROP COLUMN IF EXISTS target_summits;
