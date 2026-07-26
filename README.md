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
