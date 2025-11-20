#!/bin/sh
set -e # Salir inmediatamente si un comando falla.

echo "--- Running database migration ---"
# psql usará automáticamente la variable de entorno $DATABASE_URL,
# que es inyectada por Domcloud en este entorno de ejecución.
# El '|| true' es un seguro para que no falle si las tablas ya existen.
psql -f ./schema.sql || true
echo "--- Migration finished ---"

echo "--- Starting CollabSphere server ---"
# 'exec' reemplaza el proceso del script con el de la aplicación.
# Las variables PORT, DATABASE_URL, etc., son heredadas del entorno.
exec ./collabsphere
