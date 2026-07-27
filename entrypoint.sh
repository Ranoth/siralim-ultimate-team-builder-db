#!/bin/bash

set -e

# Load environment variables
GOOSE_DRIVER="${GOOSE_DRIVER}"
GOOSE_DBSTRING="${GOOSE_DBSTRING}"
GOOSE_MIGRATION_DIR="${GOOSE_MIGRATION_DIR}"
DB_DROP_CREATE_SEED="${DB_DROP_CREATE_SEED:-false}"

# Validate required environment variables
if [ -z "$GOOSE_DBSTRING" ]; then
    echo "Error: GOOSE_DBSTRING environment variable is not set"
    exit 1
fi

echo "Running migrations..."
echo "Driver: $GOOSE_DRIVER"
echo "Migration directory: $GOOSE_MIGRATION_DIR"

# Run goose migrations using environment variables
export GOOSE_DRIVER
export GOOSE_DBSTRING

if [ "$DB_DROP_CREATE_SEED" = "true" ]; then
	echo "Dropping DB to force re-seeding"
    cd /app && /usr/local/bin/goose -dir "$GOOSE_MIGRATION_DIR" down || echo "Warning: goose down failed, continuing..."
fi

cd /app && /usr/local/bin/goose -dir "$GOOSE_MIGRATION_DIR" up

if [ $? -ne 0 ]; then
    echo "Error: Migration failed"
    exit 1
fi

echo "Migrations completed successfully"

# Start the application
echo "Starting application..."
exec /app/sutbdb
