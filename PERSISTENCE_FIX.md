# Data Persistence Fix - Summary

## Problem
Guild configuration was being lost when rebuilding containers after `git pull`, despite previous attempts to fix persistence.

## Root Causes Identified

1. **Bind mount vulnerability**: Using `./data:/data` (bind mount) meant the data directory was part of the git working tree and could be affected by git operations like `git clean -fdx`

2. **No .gitignore protection**: The `./data` directory wasn't in `.gitignore`, leaving it vulnerable to accidental deletion

3. **Lazy directory creation**: Docker created the data directory on-demand, potentially with incorrect permissions

4. **No explicit persistence guarantee**: The rebuild process didn't clearly communicate how data is preserved

## Changes Made

### 1. Added `/data` to .gitignore
- Protects the data directory from git operations
- Prevents accidental commits or deletions

### 2. Created docker-entrypoint.sh
- Ensures `/data` directory exists with proper permissions before starting the bot
- Logs the database path being used for debugging
- Provides clear startup diagnostics

### 3. Updated Dockerfile
- Added entrypoint script execution
- Maintains proper container initialization sequence

### 4. Changed to Named Docker Volume
**Before:**
```yaml
volumes:
  - ./data:/data  # Bind mount - vulnerable to git operations
```

**After:**
```yaml
volumes:
  - weather-bot-data:/data  # Named volume - managed by Docker

volumes:
  weather-bot-data:
    name: weather-bot-data
```

**Benefits:**
- Named volumes persist independently of the host filesystem
- Not affected by git operations or working directory changes
- Docker-recommended approach for persistent data
- Survives container rebuilds, updates, and even `docker compose down` (without `-v`)

### 5. Added Database Path Logging
- Bot now logs which database file it's opening on startup
- Helps diagnose persistence issues

### 6. Updated README Documentation
- Clear instructions for rebuilding while preserving data
- Warning about `docker compose down -v` (deletes data)
- Volume management commands
- Explicit note that data survives rebuilds

## How to Test

### Fresh Installation
```bash
# Clone/pull the repo
git pull

# Start the bot (creates volume automatically)
docker compose up -d

# Configure your server
# Use /setup commands in Discord

# Verify data exists
docker compose exec bot ls -la /data
docker compose exec bot cat /data/weather-bot.db | head -c 100
```

### Rebuild After Update
```bash
# Pull latest code
git pull

# Rebuild and restart (data preserved!)
docker compose up -d --build

# Verify settings still exist
# Check Discord - your /setup configuration should still be there
```

### Verify Volume Persistence
```bash
# Check volume exists
docker volume ls | grep weather-bot-data

# Inspect volume details
docker volume inspect weather-bot-data

# View logs to confirm database path
docker compose logs | grep "DATABASE_PATH"
```

## Important Notes

### ⚠️ DO NOT use `-v` flag when stopping
```bash
# ❌ DELETES ALL DATA
docker compose down -v

# ✅ Preserves data
docker compose down
```

### Volume Management
```bash
# List all volumes
docker volume ls

# Inspect the data volume
docker volume inspect weather-bot-data

# Backup the database
docker run --rm -v weather-bot-data:/data -v $(pwd):/backup alpine tar czf /backup/db-backup.tar.gz /data

# Restore from backup
docker run --rm -v weather-bot-data:/data -v $(pwd):/backup alpine tar xzf /backup/db-backup.tar.gz -C /
```

## Migration from Old Setup

If you were using the old `./data` bind mount:

```bash
# 1. Stop the container
docker compose down

# 2. Backup existing data (if any)
cp -r ./data ./data-backup

# 3. Pull latest changes
git pull

# 4. Start with new named volume
docker compose up -d

# 5. If you had data to migrate, copy it into the volume
docker run --rm -v weather-bot-data:/data -v $(pwd)/data-backup:/backup alpine cp /backup/weather-bot.db /data/

# 6. Restart to use the migrated data
docker compose restart
```

## Verification Checklist

After rebuilding, verify:
- [ ] `docker volume ls` shows `weather-bot-data`
- [ ] Bot logs show "opening database" with path `/data/weather-bot.db`
- [ ] Discord `/setup` command shows your previous configuration
- [ ] Alert channels, language settings, and tide stations are preserved
