FROM postgres:16-alpine

COPY scripts/init.sh ./docker-entrypoint-initdb.d/
COPY sql/ ./sql
