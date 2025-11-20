#!/bin/sh

# Script de arranque para la aplicación CollabSphere en Domcloud

# 1. Ejecutar la migración de la base de datos.
#    - El comando 'psql' usará automáticamente la variable de entorno $DATABASE_URL,
#      que es inyectada por el addon de Domcloud en este entorno de ejecución.
#    - Redirigimos la salida de error a /dev/null para ignorar los errores
#      si las tablas ya existen en despliegues posteriores.
echo "--- Running database migration ---"
psql -f ./schema.sql 2>/dev/null
echo "--- Migration finished ---"

# 2. Iniciar la aplicación principal de Go.
#    - 'env PORT=$PORT' asegura que la aplicación escuche en el puerto correcto.
#    - './collabsphere' es el binario que compilamos durante la fase de 'build'.
#    - 'exec' reemplaza el proceso del script con el proceso de la aplicación,
#      lo cual es una buena práctica para la gestión de señales.
echo "--- Starting CollabSphere server ---"
exec env PORT=$PORT ./collabsphere
