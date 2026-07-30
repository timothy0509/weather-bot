# Hong Kong Weather Discord Bot

A Discord bot for Hong Kong weather, powered by the HKO (Hong Kong Observatory) Open Data API.

## Features

- Current weather, forecasts, warnings, rainfall, UV index, tides, lunar calendar, and earthquake reports
- Multi-language support: English, Traditional Chinese, Simplified Chinese, and bilingual side-by-side
- Background monitoring of HKO warnings and special weather tips
- Guild-specific settings via `/setup` commands
- SQLite persistence

## Quick Start

1. Copy `.env.example` to `.env` and fill in your Discord token.
2. Run the bot:

```bash
go run ./cmd/bot
```

Or use Docker:

```bash
docker compose up -d
```

## Docker Deployment

### Prerequisites

- Docker and Docker Compose installed on your system

### Configuration

1. Copy the environment template:

```bash
cp .env.example .env
```

2. Edit `.env` and fill in your Discord credentials:

```env
DISCORD_TOKEN=your_bot_token_here
DISCORD_APPLICATION_ID=your_application_id_here
```

Required variables:
- `DISCORD_TOKEN` — Your Discord bot token
- `DISCORD_APPLICATION_ID` — Your Discord application ID

Optional variables:
- `GUILD_ID` — Discord server ID (leave empty for global commands)
- `DATABASE_PATH` — SQLite database path (default: `./weather-bot.db`)
- `WARNING_POLL_INTERVAL` — Warning check interval (default: `90s`)
- `TIPS_POLL_INTERVAL` — Tips check interval (default: `5m`)
- `STATUS_POLL_INTERVAL` — Status update interval (default: `10m`)
- `LOG_LEVEL` — Logging level: `debug`, `info`, `warn`, `error` (default: `info`)

### Using Docker Compose (Recommended)

Build and start the bot:

```bash
docker compose up -d
```

The compose file:
- Builds the image from the local Dockerfile
- Mounts `./data` directory for SQLite database persistence
- Loads environment variables from `.env`
- Automatically restarts the container unless manually stopped

View logs:

```bash
docker compose logs -f
```

Update to latest code:

```bash
docker compose up -d --build
```

### Using Docker Only

Build the image:

```bash
docker build -t weather-bot .
```

Run the container:

```bash
docker run -d \
  --name weather-bot \
  --env-file .env \
  -v $(pwd)/data:/data \
  --restart unless-stopped \
  weather-bot
```

View logs:

```bash
docker logs -f weather-bot
```

### Data Persistence

The bot stores its SQLite database in the `/data` directory inside the container. The Docker Compose setup mounts this to `./data` on your host machine, ensuring your data persists across container restarts and updates.

**Note:** Inside the container, `DATABASE_PATH` is automatically set to `/data/weather-bot.db`.

## Commands

### Weather
- `/weather current` — Current weather
- `/weather forecast` — 3-day weather forecast
- `/weather forecast9` — 9-day weather forecast
- `/weather warnings` — Active warnings summary
- `/weather warning <type>` — Detailed warning info
- `/weather rain` — Hourly rainfall
- `/weather uv` — Current UV index
- `/weather tide [station]` — Tides at a station (default QUB)
- `/weather lunar` — Lunar calendar
- `/weather earthquake` — Earthquake reports

### Setup (Manage Server permission required)
- `/setup` — Show current settings
- `/setup language <en|tc|sc|bilingual>` — Set language mode
- `/setup alerts #channel` — Set alerts channel
- `/setup alerts disable` — Disable alerts
- `/setup tide-station <code>` — Set default tide station
- `/setup status <on|off>` — Toggle bot activity status

## Environment Variables

See `.env.example`.

## License

MIT
