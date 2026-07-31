#!/bin/sh
set -e

# Ensure data directory exists with proper permissions
if [ ! -d "/data" ]; then
    mkdir -p /data
    chmod 755 /data
fi

# Log database path for debugging
echo "Starting weather-bot with DATABASE_PATH: ${DATABASE_PATH:-/data/weather-bot.db}"

# Execute the main application
exec "$@"
