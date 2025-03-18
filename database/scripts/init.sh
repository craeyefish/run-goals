#!/bin/sh

ENVIRONMENT=$ENVIRONMENT
PGDATABASE=$POSTGRES_DB
PGUSER=$POSTGRES_USER
PGPASSWORD=$POSTGRES_PASSWORD

# DATABASE
DATABASE_DIR=../sql/db
for f in $DATABASE_DIR/*.sql;
do
    echo "$f"
    psql -U $PGUSER -f "$f"
done

# TABLES
TABLE_DIR=../sql/tables
psql -U $PGUSER -d $PGDATABASE -f "../sql/tables/users.sql"
psql -U $PGUSER -d $PGDATABASE -f "../sql/tables/peaks.sql"
psql -U $PGUSER -d $PGDATABASE -f "../sql/tables/activity.sql"
psql -U $PGUSER -d $PGDATABASE -f "../sql/tables/userpeaks.sql"

# LOAD MOCK DATA INTO DATABASE TABLES
if [ "$ENVIRONMENT" != "production" ]; then
    MOCK_DATA_DIR=../sql/mockdata
    psql -U $PGUSER -d $PGDATABASE -f "../sql/mockdata/users.sql"
    psql -U $PGUSER -d $PGDATABASE -f "../sql/mockdata/peaks.sql"
    psql -U $PGUSER -d $PGDATABASE -f "../sql/mockdata/activity.sql"
    psql -U $PGUSER -d $PGDATABASE -f "../sql/mockdata/userpeaks.sql"
fi
