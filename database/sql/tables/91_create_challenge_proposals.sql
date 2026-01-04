-- Create challenge_proposals table for community challenge ideas

CREATE TABLE IF NOT EXISTS challenge_proposals (
    id BIGSERIAL PRIMARY KEY,
    proposed_by_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Challenge details
    name VARCHAR(255) NOT NULL,
    description TEXT,
    goal_type VARCHAR(20) NOT NULL,
    competition_mode VARCHAR(20) NOT NULL DEFAULT 'competitive',

    -- Target values
    target_value NUMERIC,           -- For distance/elevation (in meters)
    target_summit_count INTEGER,    -- For summit_count
    peak_ids BIGINT[],              -- For specific_summits (array of peak IDs)

    -- Optional metadata
    region VARCHAR(100),
    difficulty VARCHAR(20),

    -- Proposal workflow
    status VARCHAR(20) DEFAULT 'pending',
    admin_notes TEXT,

    -- Timestamps
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    reviewed_at TIMESTAMPTZ,
    reviewed_by_user_id BIGINT REFERENCES users(id),

    -- Constraints
    CONSTRAINT check_proposal_status CHECK (status IN ('pending', 'approved', 'rejected')),
    CONSTRAINT check_proposal_goal_type CHECK (goal_type IN ('distance', 'elevation', 'summit_count', 'specific_summits')),
    CONSTRAINT check_proposal_competition_mode CHECK (competition_mode IN ('collaborative', 'competitive'))
);

-- Indexes for queries
CREATE INDEX idx_proposals_status ON challenge_proposals(status);
CREATE INDEX idx_proposals_user ON challenge_proposals(proposed_by_user_id);
CREATE INDEX idx_proposals_created ON challenge_proposals(created_at DESC);

-- When a proposal is approved, it will create a challenge with:
-- - visibility = 'public'
-- - is_featured = TRUE
-- - All other fields copied from proposal
