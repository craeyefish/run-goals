-- Migration: 002_create_challenges_tables.sql
-- Description: Create challenge system tables for Summit Seekers
-- Date: 2024-12-30

-- Core challenge definition
CREATE TABLE IF NOT EXISTS challenges (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    
    -- Type determines behavior
    -- 'predefined' = Admin-created, discoverable (e.g., "Cape Town 13 Peaks")
    -- 'custom'     = User-created, can be private or public
    -- 'yearly_goal'= Special type, deadline = end of year, tracks count not specific peaks
    challenge_type VARCHAR(50) NOT NULL DEFAULT 'custom',
    
    -- Competition mode
    -- 'collaborative' = Work together to complete all peaks
    -- 'competitive'   = Leaderboard, first to complete wins / most summits wins
    competition_mode VARCHAR(50) NOT NULL DEFAULT 'collaborative',
    
    -- Visibility
    -- 'private' = Only creator and invited participants
    -- 'friends' = Friends of creator can see (future feature)
    -- 'public'  = Discoverable by anyone
    visibility VARCHAR(50) NOT NULL DEFAULT 'private',
    
    -- Dates
    start_date DATE,
    deadline DATE,
    
    -- Ownership
    created_by_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    created_by_group_id BIGINT REFERENCES groups(id) ON DELETE SET NULL,
    
    -- For yearly goals specifically (target count of summits)
    target_count INT,
    
    -- Metadata for discovery/filtering
    region VARCHAR(255),
    difficulty VARCHAR(50),
    is_featured BOOLEAN DEFAULT FALSE,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Peaks in a challenge (not used for yearly_goal type)
CREATE TABLE IF NOT EXISTS challenge_peaks (
    id BIGSERIAL PRIMARY KEY,
    challenge_id BIGINT NOT NULL REFERENCES challenges(id) ON DELETE CASCADE,
    peak_id BIGINT NOT NULL REFERENCES peaks(id) ON DELETE CASCADE,
    sort_order INT DEFAULT 0,
    
    created_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(challenge_id, peak_id)
);

-- User participation in challenges
CREATE TABLE IF NOT EXISTS challenge_participants (
    id BIGSERIAL PRIMARY KEY,
    challenge_id BIGINT NOT NULL REFERENCES challenges(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    joined_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP,
    
    -- Cached progress (updated on summit detection)
    peaks_completed INT DEFAULT 0,
    total_peaks INT DEFAULT 0,
    
    UNIQUE(challenge_id, user_id)
);

-- Group participation (group adopts a challenge)
CREATE TABLE IF NOT EXISTS challenge_groups (
    id BIGSERIAL PRIMARY KEY,
    challenge_id BIGINT NOT NULL REFERENCES challenges(id) ON DELETE CASCADE,
    group_id BIGINT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    
    started_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP,
    
    -- Group can override deadline
    deadline_override DATE,
    
    UNIQUE(challenge_id, group_id)
);

-- Track which summits count toward which challenges
CREATE TABLE IF NOT EXISTS challenge_summit_log (
    id BIGSERIAL PRIMARY KEY,
    challenge_id BIGINT NOT NULL REFERENCES challenges(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    peak_id BIGINT REFERENCES peaks(id) ON DELETE SET NULL,
    activity_id BIGINT REFERENCES activity(id) ON DELETE SET NULL,
    
    summited_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    
    -- One credit per peak per user per challenge
    UNIQUE(challenge_id, user_id, peak_id)
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_challenges_type ON challenges(challenge_type);
CREATE INDEX IF NOT EXISTS idx_challenges_visibility ON challenges(visibility);
CREATE INDEX IF NOT EXISTS idx_challenges_region ON challenges(region);
CREATE INDEX IF NOT EXISTS idx_challenges_featured ON challenges(is_featured) WHERE is_featured = TRUE;
CREATE INDEX IF NOT EXISTS idx_challenges_created_by_user ON challenges(created_by_user_id);
CREATE INDEX IF NOT EXISTS idx_challenges_created_by_group ON challenges(created_by_group_id);

CREATE INDEX IF NOT EXISTS idx_challenge_peaks_challenge ON challenge_peaks(challenge_id);
CREATE INDEX IF NOT EXISTS idx_challenge_peaks_peak ON challenge_peaks(peak_id);

CREATE INDEX IF NOT EXISTS idx_challenge_participants_challenge ON challenge_participants(challenge_id);
CREATE INDEX IF NOT EXISTS idx_challenge_participants_user ON challenge_participants(user_id);
CREATE INDEX IF NOT EXISTS idx_challenge_participants_completed ON challenge_participants(completed_at) WHERE completed_at IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_challenge_groups_challenge ON challenge_groups(challenge_id);
CREATE INDEX IF NOT EXISTS idx_challenge_groups_group ON challenge_groups(group_id);

CREATE INDEX IF NOT EXISTS idx_challenge_summit_log_challenge ON challenge_summit_log(challenge_id);
CREATE INDEX IF NOT EXISTS idx_challenge_summit_log_user ON challenge_summit_log(user_id);
CREATE INDEX IF NOT EXISTS idx_challenge_summit_log_challenge_user ON challenge_summit_log(challenge_id, user_id);
