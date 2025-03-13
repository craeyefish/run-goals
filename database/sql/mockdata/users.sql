INSERT INTO
    users (
        strava_athlete_id,
        access_token,
        refresh_token,
        expires_at,
        last_distance,
        last_updated,
        created_at,
        updated_at
    )
VALUES
    (
        1001,
        'token1',
        'refresh1',
        NOW () + INTERVAL '1 hour',
        5.2,
        NOW (),
        NOW (),
        NOW ()
    ),
    (
        1002,
        'token2',
        'refresh2',
        NOW () + INTERVAL '1 hour',
        3.7,
        NOW (),
        NOW (),
        NOW ()
    ),
    (
        1003,
        'token3',
        'refresh3',
        NOW () + INTERVAL '1 hour',
        6.1,
        NOW (),
        NOW (),
        NOW ()
    ),
    (
        1004,
        'token4',
        'refresh4',
        NOW () + INTERVAL '1 hour',
        4.9,
        NOW (),
        NOW (),
        NOW ()
    ),
    (
        1005,
        'token5',
        'refresh5',
        NOW () + INTERVAL '1 hour',
        7.3,
        NOW (),
        NOW (),
        NOW ()
    );
