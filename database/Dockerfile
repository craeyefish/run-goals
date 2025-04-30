FROM postgres:16-alpine

# ---- schema ---------------------------------------------------
COPY sql/db/*.sql       /docker-entrypoint-initdb.d/10_db_/
COPY sql/tables/*.sql   /docker-entrypoint-initdb.d/20_tables_/

# ---- optional mock data --------------------------------------
ARG ENVIRONMENT=development
ENV ENVIRONMENT=${ENVIRONMENT}

# copy mock data only for non-prod images
RUN if [ "$ENVIRONMENT" != "production" ]; then \
      mkdir -p /docker-entrypoint-initdb.d/30_mockdata && \
      cp -v sql/mockdata/*.sql /docker-entrypoint-initdb.d/30_mockdata/ ; \
    fi