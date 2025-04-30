FROM postgres:16-alpine

# --- schema ----------------------------------------------------
COPY sql/db/*.sql       /docker-entrypoint-initdb.d/10_db_/
COPY sql/tables/*.sql   /docker-entrypoint-initdb.d/20_tables_/

# --- mock data -------------------------------------------------
ARG ENVIRONMENT=development
COPY sql/mockdata /tmp/mockdata

RUN if [ "$ENVIRONMENT" != "production" ]; then \
      mkdir -p /docker-entrypoint-initdb.d/30_mockdata && \
      cp -v /tmp/mockdata/*.sql /docker-entrypoint-initdb.d/30_mockdata/ ; \
    fi
