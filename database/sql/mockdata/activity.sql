INSERT INTO
    activity (
        strava_activity_id,
        strava_athlete_id,
        user_id,
        name,
        distance,
        description,
        start_date,
        map_polyline,
        created_at,
        updated_at,
        has_summit,
        photo_url
    )
VALUES
    (
        1,
        1001,
        1,
        'Morning Run',
        5.2,
        'description of the activity',
        NOW () - INTERVAL '1 day',
        'xyz123',
        NOW (),
        NOW (),
        TRUE,
        ''
    ),
    (
        2,
        1001,
        1,
        'Evening Cycle',
        15.0,
        'description of the activity',
        NOW () - INTERVAL '2 days',
        'abc456',
        NOW (),
        NOW (),
        FALSE,
        ''
    ),
    (
        3,
        1002,
        2,
        'Hiking Trip',
        7.8,
        'description of the activity',
        NOW () - INTERVAL '3 days',
        'def789',
        NOW (),
        NOW (),
        TRUE,
        ''
    ),
    (
        4,
        1002,
        2,
        'Quick Walk',
        3.2,
        'description of the activity',
        NOW () - INTERVAL '4 days',
        'ghi012',
        NOW (),
        NOW (),
        FALSE,
        ''
    ),
    (
        5,
        1003,
        3,
        'Trail Run',
        9.5,
        'description of the activity',
        NOW () - INTERVAL '5 days',
        'jkl345',
        NOW (),
        NOW (),
        TRUE,
        ''
    ),
    (
        6,
        1003,
        3,
        'City Ride',
        12.3,
        'description of the activity',
        NOW () - INTERVAL '6 days',
        'mno678',
        NOW (),
        NOW (),
        FALSE,
        ''
    ),
    (
        7,
        1004,
        4,
        'Summit Hike',
        10.0,
        'description of the activity',
        NOW () - INTERVAL '7 days',
        'pqr901',
        NOW (),
        NOW (),
        TRUE,
        ''
    ),
    (
        8,
        1004,
        4,
        'Park Jog',
        4.1,
        'description of the activity',
        NOW () - INTERVAL '8 days',
        'stu234',
        NOW (),
        NOW (),
        FALSE,
        ''
    ),
    (
        9,
        1005,
        5,
        'Long Run',
        14.2,
        'description of the activity',
        NOW () - INTERVAL '9 days',
        'vwx567',
        NOW (),
        NOW (),
        TRUE,
        ''
    ),
    (
        10,
        1005,
        5,
        'Mountain Biking',
        20.5,
        'description of the activity',
        NOW () - INTERVAL '10 days',
        'yz0123',
        NOW (),
        NOW (),
        FALSE,
        ''
    );
