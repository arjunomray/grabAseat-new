#!/bin/sh
set -e

# Apply database migrations before starting the service.
DSN="host=${DB_HOST:-localhost} port=${DB_PORT:-5432} user=${DB_USER:-postgres} password=${DB_PASSWORD:-postgres} dbname=${DB_NAME:-eventer} sslmode=disable"
goose -dir /app/migrations postgres "$DSN" up

exec /app/eventer
