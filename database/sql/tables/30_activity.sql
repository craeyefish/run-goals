CREATE TABLE IF NOT EXISTS activity (
    id BIGSERIAL PRIMARY KEY,
    strava_activity_id BIGINT,
    strava_athlete_id BIGINT,
    user_id BIGINT,
    name VARCHAR,
    description VARCHAR,
    distance NUMERIC,
    start_date TIMESTAMPTZ,
    map_polyline VARCHAR,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    has_summit BOOLEAN,
    photo_url VARCHAR,
    CONSTRAINT fk_activity_user_activity FOREIGN KEY (strava_athlete_id) REFERENCES users (strava_athlete_id) ON DELETE CASCADE,
    CONSTRAINT activity_strava_activity_id_key UNIQUE (strava_activity_id)
);
