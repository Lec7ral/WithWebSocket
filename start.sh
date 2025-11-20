#!/bin/sh
set -e # Exit immediately if a command exits with a non-zero status.

echo "--- Running database migration ---"
# El addon de Domcloud inyecta $DATABASE_URL en este entorno.
# psql lo usará automáticamente para conectarse y crear las tablas.
psql -f ./schema.sql
echo "--- Migration finished ---"

echo "--- Starting CollabSphere server ---"
# Le pasamos la DATABASE_URL a nuestra aplicación como una variable de entorno
# para que Viper la lea.
exec env PORT=$PORT DATABASE_URL=$DATABASE_URL ./collabsphere
