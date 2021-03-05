#!/bin/bash
set -e

psql \c "$POSTGRES_DB" "$POSTGRES_USER" <<-EOSQL
    CREATE TABLE IF NOT EXISTS Users (ID serial PRIMARY KEY, Name VARCHAR (50) NOT NULL, Country VARCHAR (10)  NOT NULL, Points INT NOT NULL);    
EOSQL
