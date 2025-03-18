CREATE TABLE IF NOT EXISTS group_goals (
    id BIGSERIAL PRIMARY KEY,
    group_id BIGINT,
    name VARCHAR,
    target_value NUMERIC,
    start_date TIMESTAMPZ,
    end_date TIMESTAMPZ,
    created_at TIMESTAMPZ,
    CONSTRAINT fk_group_goals_group FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE
);
