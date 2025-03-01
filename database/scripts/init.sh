#!/bin/sh

PGHOST="localhost"
PGDATABASE="rungoals"
PGUSER="postgres"
PGPASSWORD="postgres"

# DATABASE
DATABASE_DIR=../sql/db
for f in $DATABASE_DIR/*.sql;
do
    echo "$f"
    psql -U $PGUSER -f "$f"
done

# TABLES
TABLE_DIR=../sql/tables
psql -U $PGUSER -d $PGDATABASE -f "../sql/tables/user.sql"
psql -U $PGUSER -d $PGDATABASE -f "../sql/tables/peak.sql"
psql -U $PGUSER -d $PGDATABASE -f "../sql/tables/activity.sql"
psql -U $PGUSER -d $PGDATABASE -f "../sql/tables/userpeak.sql"

# LOAD MOCK DATA INTO DATABASE TABLES
# MOCK_DATA_DIR=../sql/data
# psql -U $PGUSER -d $PGDATABASE -f "../sql/data/<file_name>.sql"
