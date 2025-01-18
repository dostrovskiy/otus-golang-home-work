#!/bin/sh
sleep 5

echo "Creating database..."
PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -U "$DB_USER" -d postgres -c "CREATE DATABASE $DB_NAME"
echo "Database $DB_NAME created"

sleep 5

# Выполнение миграций
GOOSE_DRIVER=postgres GOOSE_DBSTRING=${DB_STRING} GOOSE_MIGRATION_DIR=/migrations goose up

sleep 5