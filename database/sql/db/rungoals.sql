SELECT
  'CREATE DATABASE rungoals'
WHERE NOT EXISTS (
  SELECT FROM pg_database WHERE datname = 'rungoals'
); \gexec
