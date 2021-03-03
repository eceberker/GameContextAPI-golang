#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER postgres;
    CREATE DATABASE GameContextDb;
    GRANT ALL PRIVILEGES ON DATABASE GameContextDb TO postgres;
    CREATE DATABASE GameContextDb;
    GRANT ALL PRIVILEGES ON DATABASE GameContextDb TO postgres;

    CREATE TABLE IF NOT EXISTS Users (ID serial PRIMARY KEY, Name VARCHAR (50) NOT NULL, Country VARCHAR (10)  NOT NULL, Points INT NOT NULL);    
EOSQL