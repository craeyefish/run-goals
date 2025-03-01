CREATE TABLE IF NOT EXISTS user_peak
(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    peak_id BIGINT,
    activity_id BIGINT,
    summited_at TIMESTAMPTZ,
    CONSTRAINT fk_user_peak_user FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE RESTRICT,
    CONSTRAINT fk_user_peak_peak_id FOREIGN KEY (peak_id) REFERENCES peak(id) ON DELETE RESTRICT,
    CONSTRAINT fk_user_peak_activity_id FOREIGN KEY (activity_id) REFERENCES activity(id) ON DELETE RESTRICT
);
