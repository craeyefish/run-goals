CREATE TABLE IF NOT EXISTS activity (
    id BIGSERIAL PRIMARY KEY,
    strava_athlete_id BIGINT,
    user_id BIGINT,
    name VARCHAR,
    distance NUMERIC,
    start_date TIMESTAMPTZ,
    map_polyline VARCHAR,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    has_summit BOOLEAN,
    CONSTRAINT fk_activity_user_activity FOREIGN KEY (strava_athlete_id) REFERENCES users (strava_athlete_id) ON DELETE CASCADE,
    CONSTRAINT fk_activity_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
