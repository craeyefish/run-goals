-- Personal yearly goals for tracking individual user targets per year
CREATE TABLE IF NOT EXISTS personal_yearly_goals (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    year INT NOT NULL,
    distance_goal NUMERIC DEFAULT 0,           -- Target distance in km
    elevation_goal NUMERIC DEFAULT 0,          -- Target elevation in meters
    summit_goal INT DEFAULT 0,                 -- Number of summits to achieve
    target_summits BIGINT[],                   -- Specific peak IDs user wants to summit
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_personal_goal_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT unique_user_year UNIQUE (user_id, year),
    CONSTRAINT check_year CHECK (year >= 2020 AND year <= 2100),
    CONSTRAINT check_distance_goal CHECK (distance_goal >= 0),
    CONSTRAINT check_elevation_goal CHECK (elevation_goal >= 0),
    CONSTRAINT check_summit_goal CHECK (summit_goal >= 0)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_personal_yearly_goals_user_id ON personal_yearly_goals(user_id);
CREATE INDEX IF NOT EXISTS idx_personal_yearly_goals_year ON personal_yearly_goals(year);
CREATE INDEX IF NOT EXISTS idx_personal_yearly_goals_user_year ON personal_yearly_goals(user_id, year);
