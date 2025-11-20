#!/bin/sh
set -e # Exit immediately if a command exits with a non-zero status.

echo "--- Running database migration ---"
# psql usará la DATABASE_URL que definimos en domcloud.yml
# y que se inyecta en este entorno de ejecución.
psql -f ./schema.sql
echo "--- Migration finished ---"

echo "--- Starting CollabSphere server ---"
exec env PORT=$PORT ./collabsphere
