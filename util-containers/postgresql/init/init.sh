#!/bin/bash

# This tells bash that it should exit the script if any statement returns a non-true return value.
# http://web.archive.org/web/20110314180918/http://www.davidpashley.com/articles/writing-robust-shell-scripts.html
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE "nulab_backlog" ENCODING 'UTF8' TEMPLATE template0;
EOSQL
