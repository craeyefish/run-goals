FROM postgres:16-alpine

RUN echo '\connect run_goals' > /docker-entrypoint-initdb.d/00_connect.sql

# ---------------- schema (always) ----------------
COPY sql/db/*.sql       /docker-entrypoint-initdb.d/
COPY sql/tables/*.sql   /docker-entrypoint-initdb.d/

# ---------------- mock data (non-prod) -----------
  # # Uncomment the following lines to enable mock data loading
  
# ARG ENVIRONMENT=development
# # copy into a temp dir first
# COPY sql/mockdata /tmp/mockdata
# RUN if [ "$ENVIRONMENT" != "production" ]; then \
#       cp /tmp/mockdata/*.sql /docker-entrypoint-initdb.d/ ; \
#     fi