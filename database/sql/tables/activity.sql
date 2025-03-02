CREATE TABLE IF NOT EXISTS activity
(
    id BIGSERIAL PRIMARY KEY,
    strava_athlete_id BIGINT CONSTRAINT unique_activity_strava_athlete_id UNIQUE,
    user_id BIGINT,
    name VARCHAR,
    distance NUMERIC,
    start_date TIMESTAMPTZ,
    map_polyline VARCHAR,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    has_summit BOOLEAN,
    CONSTRAINT fk_activity_user FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE
);
