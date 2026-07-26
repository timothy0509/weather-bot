# Agent Notes

## Project

- Go 1.26+ Discord bot for HK weather using HKO Open Data API
- Pure Go SQLite via `modernc.org/sqlite`
- Uses `github.com/bwmarrin/discordgo`

## Build & Run

```bash
go run ./cmd/bot
```

## Test

```bash
go test ./...
```

## Key Files

- `cmd/bot/main.go` — Entry point
- `internal/config/` — Configuration
- `internal/db/` — SQLite persistence
- `internal/hko/` — HKO API client
- `internal/i18n/` — Multi-language strings
- `internal/bot/` — Discord bot logic
- `internal/monitor/` — Background monitoring

## Conventions

- English variable/comment names, i18n for user-facing text
- Minimal changes to existing code; follow existing style
- `/setup` commands require `Manage Server` permission
- Default tide station is QUB
- Language modes: en, tc, sc, bilingual
