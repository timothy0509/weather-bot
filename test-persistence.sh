#!/bin/bash
# Test script to verify data persistence is working correctly

set -e

echo "🔍 Testing Data Persistence Setup..."
echo ""

# Check if docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

echo "✓ Docker is running"

# Check if the volume exists
if docker volume inspect weather-bot-data > /dev/null 2>&1; then
    echo "✓ Named volume 'weather-bot-data' exists"
    
    # Check volume details
    MOUNTPOINT=$(docker volume inspect weather-bot-data --format '{{ .Mountpoint }}')
    echo "  Volume mountpoint: $MOUNTPOINT"
else
    echo "⚠ Volume 'weather-bot-data' does not exist yet"
    echo "  (It will be created when you run 'docker compose up -d')"
fi

# Check if container is running
if docker ps --format '{{.Names}}' | grep -q "^weather-bot$"; then
    echo "✓ Container 'weather-bot' is running"
    
    # Check if database file exists in container
    if docker exec weather-bot test -f /data/weather-bot.db 2>/dev/null; then
        echo "✓ Database file exists at /data/weather-bot.db"
        
        # Get database size
        DB_SIZE=$(docker exec weather-bot ls -lh /data/weather-bot.db | awk '{print $5}')
        echo "  Database size: $DB_SIZE"
    else
        echo "⚠ Database file not found in container"
    fi
    
    # Check recent logs for database path
    echo ""
    echo "📋 Recent startup logs:"
    docker logs weather-bot --tail 20 2>&1 | grep -E "(DATABASE_PATH|opening database)" || echo "  (No recent startup logs found)"
else
    echo "⚠ Container 'weather-bot' is not running"
    echo "  Run 'docker compose up -d' to start it"
fi

# Check .gitignore
echo ""
if grep -q "^/data$" .gitignore 2>/dev/null; then
    echo "✓ /data is in .gitignore"
else
    echo "⚠ /data is not in .gitignore"
fi

# Check entrypoint script
if [ -f "docker-entrypoint.sh" ]; then
    echo "✓ docker-entrypoint.sh exists"
    if [ -x "docker-entrypoint.sh" ]; then
        echo "✓ docker-entrypoint.sh is executable"
    else
        echo "⚠ docker-entrypoint.sh is not executable"
    fi
else
    echo "❌ docker-entrypoint.sh not found"
fi

# Check docker-compose.yml
if grep -q "weather-bot-data:/data" docker-compose.yml 2>/dev/null; then
    echo "✓ docker-compose.yml uses named volume"
else
    echo "⚠ docker-compose.yml may not be using named volume"
fi

echo ""
echo "✅ Persistence setup verification complete!"
echo ""
echo "📝 To test persistence:"
echo "   1. Configure your bot using /setup commands in Discord"
echo "   2. Run: git pull"
echo "   3. Run: docker compose up -d --build"
echo "   4. Verify your settings are still there with /setup"
