#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER postgres;
    CREATE DATABASE GameContextDb;
    GRANT ALL PRIVILEGES ON DATABASE GameContextDb TO postgres;
    CREATE DATABASE GameContextDb;
    GRANT ALL PRIVILEGES ON DATABASE GameContextDb TO postgres;
EOSQL