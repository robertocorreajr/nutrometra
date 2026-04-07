#!/bin/bash
# Creates the application database, Zitadel database, and app user.
# Reads credentials from environment variables set in docker-compose.yml.
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Application database
    SELECT 'CREATE DATABASE nutrometra'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'nutrometra')\gexec

    -- Zitadel database
    SELECT 'CREATE DATABASE zitadel'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'zitadel')\gexec

    -- Application user
    DO \$\$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = '${APP_DB_USER}') THEN
            EXECUTE format('CREATE USER %I WITH PASSWORD %L', '${APP_DB_USER}', '${APP_DB_PASSWORD}');
        END IF;
    END
    \$\$;

    GRANT ALL PRIVILEGES ON DATABASE nutrometra TO ${APP_DB_USER};
EOSQL

# PostgreSQL 15+ requires explicit schema permission
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "nutrometra" <<-EOSQL
    GRANT ALL ON SCHEMA public TO ${APP_DB_USER};
EOSQL
