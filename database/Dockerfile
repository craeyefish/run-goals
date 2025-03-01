FROM postgres:16-alpine

ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD postgres
ENV POSTGRES_DB rungoals

COPY scripts/init.sh ./docker-entrypoint-initdb.d/
COPY sql/ ./sql
