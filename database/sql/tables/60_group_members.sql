CREATE TABLE IF NOT EXISTS group_members (
    id BIGSERIAL PRIMARY KEY,
    group_id BIGINT,
    user_id BIGINT,
    role VARCHAR,
    joined_at TIMESTAMPTZ,
    CONSTRAINT fk_group_members_group FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE,
    CONSTRAINT fk_group_members_users FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
