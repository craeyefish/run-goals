CREATE TABLE IF NOT EXISTS "user"
(
    id BIGSERIAL PRIMARY KEY,
    strava_athlete_id BIGINT CONSTRAINT unique_user_strava_athlete_id UNIQUE,
    access_token VARCHAR,
    refresh_token VARCHAR,
    expires_at TIMESTAMPTZ,
    last_distance NUMERIC,
    last_updated TIMESTAMPTZ,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);
