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
- Uses a named Docker volume (`weather-bot-data`) for SQLite database persistence
- Loads environment variables from `.env`
- Automatically restarts the container unless manually stopped

View logs:

```bash
docker compose logs -f
```

Update to latest code (data is preserved):

```bash
git pull
docker compose up -d --build
```

**Note:** Your guild settings and data are stored in a Docker named volume and will persist across container rebuilds. The volume is created automatically on first run.

### Using Docker Only

Build the image:

```bash
docker build -t weather-bot .
```

Run the container with a named volume:

```bash
docker volume create weather-bot-data
docker run -d \
  --name weather-bot \
  --env-file .env \
  -v weather-bot-data:/data \
  --restart unless-stopped \
  weather-bot
```

View logs:

```bash
docker logs -f weather-bot
```

**Note:** Use `docker stop weather-bot` to stop the container (preserves data). Avoid `docker rm -v weather-bot` as the `-v` flag removes the volume and all data.

### Data Persistence

The bot stores its SQLite database in the `/data` directory inside the container. Docker Compose uses a **named volume** (`weather-bot-data`) to persist this data, ensuring your guild settings and data survive container rebuilds and updates.

**Inside the container:** `DATABASE_PATH` is automatically set to `/data/weather-bot.db`

**Managing the volume:**
- View volume: `docker volume inspect weather-bot-data`
- List volumes: `docker volume ls`
- ⚠️ **Warning:** Running `docker compose down -v` will DELETE the volume and all data. Use `docker compose down` (without `-v`) to stop containers while preserving data.

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

## Data Sources

All weather data is provided by the [Hong Kong Observatory (HKO) Open Data API](https://data.weather.gov.hk/weatherAPI/doc/HKO_Open_Data_API_Documentation.pdf). No API key is required.

| Endpoint | Description |
|---|---|
| `weather.php?dataType=rhrread` | Current weather conditions (temperature, humidity, rainfall, UV, warnings) |
| `weather.php?dataType=fnd` | 9-day weather forecast |
| `weather.php?dataType=warnsum` | Active weather warnings summary |
| `weather.php?dataType=warningInfo` | Detailed warning bulletins |
| `weather.php?dataType=swt` | Special weather tips |
| `weather.php?dataType=hourlyRainfall` | Hourly rainfall data from automatic stations |
| `earthquake.php?dataType=qem` | Recent earthquake reports |
| `opendata.php?dataType=HHOT` | Hourly astronomical tide heights (14 stations) |
| `lunardate.php` | Gregorian-to-Lunar calendar conversion |

## Environment Variables

See `.env.example`.

## License

MIT
