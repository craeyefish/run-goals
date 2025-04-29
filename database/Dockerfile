FROM postgres:16-alpine

# DO NOT run psql or copy a pre-populated /var/lib/postgresql/data
# Just drop your init scripts in the magic dir:
COPY scripts/init.sh /docker-entrypoint-initdb.d/00_init.sh
COPY sql/db/*.sql       /docker-entrypoint-initdb.d/10_schema_/
COPY sql/tables/*.sql   /docker-entrypoint-initdb.d/20_tables_/
COPY sql/mockdata/*.sql /docker-entrypoint-initdb.d/30_mockdata_/

# be sure nothing ended up in the data dir
RUN rm -rf /var/lib/postgresql/data/*