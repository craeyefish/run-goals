CREATE TABLE IF NOT EXISTS group_goals (
    id BIGSERIAL PRIMARY KEY,
    group_id BIGINT,
    name VARCHAR,
    target_value NUMERIC,
    start_date TIMESTAMPTZ,
    end_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ,
    CONSTRAINT fk_group_goals_group FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE
);
