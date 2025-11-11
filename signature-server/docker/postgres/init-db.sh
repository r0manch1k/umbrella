#!/bin/sh
set -e

psql -v ON_ERROR_STOP=0 --username "$POSTGRES_USER" --dbname=postgres <<-EOSQL
    CREATE DATABASE ${POSTGRES_DB};
EOSQL
