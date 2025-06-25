CREATE TABLE IF NOT EXISTS group_goals (
    id BIGSERIAL PRIMARY KEY,
    group_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    goal_type VARCHAR(20) NOT NULL DEFAULT 'distance',
    target_value NUMERIC NOT NULL,
    target_summits BIGINT[],
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_group_goal_group_id FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE,
    CONSTRAINT check_goal_type CHECK (goal_type IN ('distance', 'elevation', 'summit_count', 'specific_summits')),
    CONSTRAINT check_dates CHECK (end_date > start_date),
    CONSTRAINT check_target_value CHECK (target_value > 0)
);

-- Index for performance
CREATE INDEX IF NOT EXISTS idx_group_goals_group_id ON group_goals(group_id);
CREATE INDEX IF NOT EXISTS idx_group_goals_dates ON group_goals(start_date, end_date);
CREATE INDEX IF NOT EXISTS idx_group_goals_type ON group_goals(goal_type);