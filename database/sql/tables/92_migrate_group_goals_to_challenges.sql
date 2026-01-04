-- Migrate existing group_goals to challenges
-- This allows us to deprecate the group_goals table and unify on challenges

-- Step 1: Create challenges from existing group_goals
INSERT INTO challenges (
    name,
    description,
    challenge_type,
    goal_type,
    competition_mode,
    visibility,
    start_date,
    deadline,
    created_by_group_id,
    target_value,
    target_summit_count,
    region,
    difficulty,
    is_featured,
    join_code,
    is_locked,
    created_at
)
SELECT
    gg.name,
    gg.description,
    'custom'::VARCHAR,  -- Group goals are custom challenges
    gg.goal_type,
    'collaborative'::VARCHAR,  -- Group goals are always collaborative
    'private'::VARCHAR,  -- Group goals are private to the group
    gg.start_date,
    gg.end_date,
    gg.group_id,
    -- Set target_value for distance/elevation goals
    CASE
        WHEN gg.goal_type IN ('distance', 'elevation') THEN gg.target_value
        ELSE NULL
    END,
    -- Set target_summit_count for summit_count goals
    CASE
        WHEN gg.goal_type = 'summit_count' THEN gg.target_value::INTEGER
        ELSE NULL
    END,
    NULL,  -- No region in group_goals
    'medium'::VARCHAR,  -- Default difficulty
    FALSE,  -- Not featured
    UPPER(SUBSTRING(MD5(RANDOM()::TEXT || gg.id::TEXT || NOW()::TEXT) FROM 1 FOR 6)),  -- Generate unique join code
    FALSE,  -- Not locked
    gg.created_at
FROM group_goals gg
ON CONFLICT (join_code) DO NOTHING;  -- Skip if join code collision (very unlikely)

-- Step 2: Create challenge_peaks for specific_summits goals
INSERT INTO challenge_peaks (challenge_id, peak_id, sort_order)
SELECT
    c.id,
    unnest(gg.target_summits) as peak_id,
    generate_series(1, array_length(gg.target_summits, 1)) as sort_order
FROM group_goals gg
JOIN challenges c ON
    c.name = gg.name
    AND c.created_by_group_id = gg.group_id
    AND c.goal_type = 'specific_summits'
WHERE gg.goal_type = 'specific_summits'
  AND gg.target_summits IS NOT NULL
  AND array_length(gg.target_summits, 1) > 0
ON CONFLICT (challenge_id, peak_id) DO NOTHING;

-- Step 3: Auto-join all group members to migrated challenges
INSERT INTO challenge_participants (challenge_id, user_id, total_peaks, joined_at)
SELECT DISTINCT
    c.id,
    gm.user_id,
    COALESCE((SELECT COUNT(*) FROM challenge_peaks cp WHERE cp.challenge_id = c.id), 0),
    COALESCE(c.created_at, NOW())
FROM challenges c
JOIN group_members gm ON gm.group_id = c.created_by_group_id
WHERE c.created_by_group_id IS NOT NULL
ON CONFLICT (challenge_id, user_id) DO NOTHING;

-- Step 4: Link challenges to groups via challenge_groups table
INSERT INTO challenge_groups (challenge_id, group_id, started_at)
SELECT
    c.id,
    c.created_by_group_id,
    COALESCE(c.created_at, NOW())
FROM challenges c
WHERE c.created_by_group_id IS NOT NULL
ON CONFLICT (challenge_id, group_id) DO NOTHING;

-- Step 5: For audit purposes, add a comment to track migration
COMMENT ON TABLE group_goals IS 'DEPRECATED: Migrated to challenges table. Keep for rollback safety. Can be dropped in future version.';

-- Note: We are NOT dropping group_goals table yet for rollback safety
-- The frontend will stop using it, but the data remains for a grace period
-- After confirming migration success, we can drop it in a future migration
