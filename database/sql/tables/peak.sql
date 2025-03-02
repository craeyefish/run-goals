CREATE TABLE IF NOT EXISTS peak
(
    id BIGSERIAL PRIMARY KEY,
    osm_id BIGINT CONSTRAINT unique_peak_osm_id UNIQUE,
    latitude NUMERIC,
    longitude NUMERIC,
    name VARCHAR,
    elevation_meters NUMERIC
);
