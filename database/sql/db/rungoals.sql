SELECT
  'CREATE DATABASE run_goals'
WHERE NOT EXISTS (
  SELECT FROM pg_database WHERE datname = 'run_goals'
); \gexec
