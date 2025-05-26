CREATE TABLE IF NOT EXISTS user_peaks (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    peak_id BIGINT,
    activity_id BIGINT,
    summited_at TIMESTAMPTZ,
    UNIQUE (user_id, peak_id),
    CONSTRAINT fk_user_peak_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_user_peak_peak_id FOREIGN KEY (peak_id) REFERENCES peaks (id) ON DELETE CASCADE,
    CONSTRAINT fk_user_peak_activity_id FOREIGN KEY (activity_id) REFERENCES activity (id) ON DELETE CASCADE
);
