#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE db_smartfactory_iot;
    CREATE DATABASE db_smartfactory_production;
    CREATE DATABASE db_smartfactory_qc;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "db_smartfactory_production" -f /docker-entrypoint-initdb.d/schema.sql

