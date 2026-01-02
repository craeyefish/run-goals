-- Summit favourites (wishlist) - year-independent peak tracking
CREATE TABLE IF NOT EXISTS summit_favourites (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    peak_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_summit_favourite_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_summit_favourite_peak_id FOREIGN KEY (peak_id) REFERENCES peaks (id) ON DELETE CASCADE,
    CONSTRAINT unique_user_peak_favourite UNIQUE (user_id, peak_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_summit_favourites_user_id ON summit_favourites(user_id);
CREATE INDEX IF NOT EXISTS idx_summit_favourites_peak_id ON summit_favourites(peak_id);
