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
psql -U $PGUSER -d $PGDATABASE -f "../sql/tables/10_users.sql"
psql -U $PGUSER -d $PGDATABASE -f "../sql/tables/20_peaks.sql"
psql -U $PGUSER -d $PGDATABASE -f "../sql/tables/30_activity.sql"
psql -U $PGUSER -d $PGDATABASE -f "../sql/tables/40_userpeaks.sql"
psql -U $PGUSER -d $PGDATABASE -f "../sql/tables/50_groups.sql"
psql -U $PGUSER -d $PGDATABASE -f "../sql/tables/60_group_members.sql"
psql -U $PGUSER -d $PGDATABASE -f "../sql/tables/70_group_goals.sql"

# LOAD MOCK DATA INTO DATABASE TABLES
# if [ "$ENVIRONMENT" != "production" ]; then
#     MOCK_DATA_DIR=../sql/mockdata
#     psql -U $PGUSER -d $PGDATABASE -f "../sql/mockdata/users.sql"
#     psql -U $PGUSER -d $PGDATABASE -f "../sql/mockdata/peaks.sql"
#     psql -U $PGUSER -d $PGDATABASE -f "../sql/mockdata/activity.sql"
#     psql -U $PGUSER -d $PGDATABASE -f "../sql/mockdata/userpeaks.sql"
#     psql -U $PGUSER -d $PGDATABASE -f "../sql/mockdata/groups.sql"
#     psql -U $PGUSER -d $PGDATABASE -f "../sql/mockdata/group_members.sql"
#     psql -U $PGUSER -d $PGDATABASE -f "../sql/mockdata/group_goals.sql"
# fi
