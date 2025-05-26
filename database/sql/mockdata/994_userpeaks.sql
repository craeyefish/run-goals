INSERT INTO
    user_peaks (user_id, peak_id, activity_id, summited_at)
VALUES
    (1, 1, 1, NOW () - INTERVAL '1 day'),
    (2, 2, 3, NOW () - INTERVAL '3 days'),
    (3, 3, 5, NOW () - INTERVAL '5 days'),
    (4, 4, 7, NOW () - INTERVAL '7 days'),
    (5, 5, 9, NOW () - INTERVAL '9 days');
