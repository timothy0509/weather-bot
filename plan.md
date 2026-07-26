# Hong Kong Weather Discord Bot — Implementation Plan

## Stack

| Layer | Choice |
|---|---|
| Language | Go 1.22+ |
| Discord library | `github.com/bwmarrin/discordgo` |
| HTTP client | `net/http` + `encoding/json` |
| Environment | `github.com/joho/godotenv` |
| Database | `modernc.org/sqlite` (pure Go, no CGO) |
| Background tasks | `time.Ticker` + goroutines |
| Logging | `log/slog` |
| Deployment | Single binary + optional `Dockerfile`/`docker-compose.yml` |

## API Validation

| Endpoint | Status |
|---|---|
| `weather.php?dataType=rhrread&lang=en` | ✅ Live |
| `weather.php?dataType=fnd&lang=en` | ✅ Live |
| `weather.php?dataType=warnsum&lang=en` | ✅ Live (Amber Rainstorm + Thunderstorm active) |
| `opendata.php?dataType=HHOT&station=QUB` | ✅ Valid (Quarry Bay default) |

Default tide station: **QUB** (Quarry Bay).

## Language Options

| Code | Mode |
|---|---|
| `en` | English only |
| `tc` | Traditional Chinese only |
| `sc` | Simplified Chinese only |
| `bilingual` | English + Traditional Chinese side-by-side |

Bilingual embeds use side-by-side format:
```
Temperature / 氣溫
28°C / 28°C
```

## Database Schema

```sql
CREATE TABLE guild_settings (
    guild_id TEXT PRIMARY KEY,
    alert_channel_id TEXT,
    language TEXT NOT NULL DEFAULT 'en' CHECK(language IN ('en', 'tc', 'sc', 'bilingual')),
    tide_station TEXT NOT NULL DEFAULT 'QUB',
    bot_status_enabled BOOLEAN NOT NULL DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE warning_state (
    id INTEGER PRIMARY KEY,
    code TEXT NOT NULL,
    subtype TEXT,
    action_code TEXT,
    issue_time TEXT,
    update_time TEXT,
    last_seen DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tips_state (
    id INTEGER PRIMARY KEY,
    update_time TEXT NOT NULL,
    last_seen DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## Commands

### Public Weather Commands

| Command | API |
|---|---|
| `/weather current` | `rhrread` |
| `/weather forecast` | `fnd` (3 days) |
| `/weather forecast9` | `fnd` (9 days) |
| `/weather warnings` | `warnsum` |
| `/weather warning <type>` | `warningInfo` |
| `/weather rain` | `hourlyRainfall` |
| `/weather uv` | `rhrread` |
| `/weather tide [station]` | `opendata.php?HHOT` |
| `/weather lunar` | `lunardate.php` |
| `/weather earthquake` | `earthquake.php?qem` |

### Admin Setup Commands (Manage Server permission)

| Command | Description |
|---|---|
| `/setup` | Show current guild settings |
| `/setup language <en\|tc\|sc\|bilingual>` | Set language mode |
| `/setup alerts #channel` | Set alert channel |
| `/setup alerts disable` | Remove alert channel |
| `/setup tide-station <code>` | Set default tide station |
| `/setup status <on\|off>` | Toggle bot activity status |

## Background Monitoring

| Loop | Interval | Notes |
|---|---|---|
| Warning monitor | 90s | Detect new/escalated/cancelled warnings, post to alert channel |
| Special tips monitor | 5m | Post new SWT tips to alert channel |
| Status monitor | 10m | Update bot activity from `rhrread` if enabled |

If no alert channel is configured, background alerts are skipped silently.

## Project Structure

```
weather-bot/
├── .env.example
├── .gitignore
├── go.mod
├── main.go
├── Dockerfile
├── docker-compose.yml
├── plan.md
├── README.md
├── AGENTS.md
├── cmd/
│   └── bot/
│       └── main.go
└── internal/
    ├── bot/
    │   ├── bot.go
    │   ├── commands.go
    │   └── handlers.go
    ├── config/
    │   └── config.go
    ├── db/
    │   ├── db.go
    │   └── models.go
    ├── hko/
    │   ├── client.go
    │   ├── weather.go
    │   ├── warnings.go
    │   ├── earthquake.go
    │   ├── rainfall.go
    │   ├── opendata.go
    │   └── types.go
    ├── i18n/
    │   ├── i18n.go
    │   └── messages.go
    └── monitor/
        └── monitor.go
```

## Implementation Phases

1. **Foundation**: Go module, config, SQLite, HKO client.
2. **Core Bot**: Discord session, slash commands, handlers.
3. **Weather Commands**: All `/weather` commands.
4. **Admin Setup**: `/setup` commands with permission checks.
5. **Background Monitoring**: Warning/tips/status loops.
6. **Language & Polish**: Multi-language support, side-by-side embeds, logging.
7. **Testing & Deployment**: Unit tests, Dockerfile, docker-compose.

## Decisions

- Default tide station: **QUB**.
- `/setup` commands restricted to `Manage Server`.
- Background alerts skipped if no alert channel is set.
- Bot activity/status enabled by default.
- Bilingual mode uses `en` + `tc` side-by-side.

## Environment Variables

```env
DISCORD_TOKEN=
DISCORD_APPLICATION_ID=
GUILD_ID=                    # Optional, for faster guild command testing
DATABASE_PATH=./weather-bot.db
WARNING_POLL_INTERVAL=90s
TIPS_POLL_INTERVAL=5m
STATUS_POLL_INTERVAL=10m
LOG_LEVEL=info
```
