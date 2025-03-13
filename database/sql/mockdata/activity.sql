INSERT INTO
    activity (
        strava_athlete_id,
        user_id,
        name,
        distance,
        start_date,
        map_polyline,
        created_at,
        updated_at,
        has_summit
    )
VALUES
    (
        1001,
        1,
        'Morning Run',
        5.2,
        NOW () - INTERVAL '1 day',
        'xyz123',
        NOW (),
        NOW (),
        TRUE
    ),
    (
        1001,
        1,
        'Evening Cycle',
        15.0,
        NOW () - INTERVAL '2 days',
        'abc456',
        NOW (),
        NOW (),
        FALSE
    ),
    (
        1002,
        2,
        'Hiking Trip',
        7.8,
        NOW () - INTERVAL '3 days',
        'def789',
        NOW (),
        NOW (),
        TRUE
    ),
    (
        1002,
        2,
        'Quick Walk',
        3.2,
        NOW () - INTERVAL '4 days',
        'ghi012',
        NOW (),
        NOW (),
        FALSE
    ),
    (
        1003,
        3,
        'Trail Run',
        9.5,
        NOW () - INTERVAL '5 days',
        'jkl345',
        NOW (),
        NOW (),
        TRUE
    ),
    (
        1003,
        3,
        'City Ride',
        12.3,
        NOW () - INTERVAL '6 days',
        'mno678',
        NOW (),
        NOW (),
        FALSE
    ),
    (
        1004,
        4,
        'Summit Hike',
        10.0,
        NOW () - INTERVAL '7 days',
        'pqr901',
        NOW (),
        NOW (),
        TRUE
    ),
    (
        1004,
        4,
        'Park Jog',
        4.1,
        NOW () - INTERVAL '8 days',
        'stu234',
        NOW (),
        NOW (),
        FALSE
    ),
    (
        1005,
        5,
        'Long Run',
        14.2,
        NOW () - INTERVAL '9 days',
        'vwx567',
        NOW (),
        NOW (),
        TRUE
    ),
    (
        1005,
        5,
        'Mountain Biking',
        20.5,
        NOW () - INTERVAL '10 days',
        'yz0123',
        NOW (),
        NOW (),
        FALSE
    );
