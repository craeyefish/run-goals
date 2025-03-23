CREATE TABLE IF NOT EXISTS groups (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR,
    created_by BIGINT,
    created_at TIMESTAMPTZ,
    CONSTRAINT fk_groups_user FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE CASCADE
);
