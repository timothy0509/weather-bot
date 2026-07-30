package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration.
type Config struct {
	DiscordToken     string
	ApplicationID    string
	GuildID          string
	DatabasePath     string
	WarningInterval  time.Duration
	TipsInterval     time.Duration
	StatusInterval   time.Duration
	LogLevel         slog.Level
}

// Load reads configuration from the environment.
func Load(envFiles ...string) (*Config, error) {
	_ = godotenv.Load(envFiles...)

	cfg := &Config{
		DiscordToken:  os.Getenv("DISCORD_TOKEN"),
		ApplicationID: os.Getenv("DISCORD_APPLICATION_ID"),
		GuildID:       os.Getenv("GUILD_ID"),
		DatabasePath:  getEnv("DATABASE_PATH", "./weather-bot.db"),
		LogLevel:      parseLogLevel(getEnv("LOG_LEVEL", "info")),
	}

	if cfg.DiscordToken == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN is required")
	}
	if cfg.ApplicationID == "" {
		return nil, fmt.Errorf("DISCORD_APPLICATION_ID is required")
	}
	if cfg.GuildID != "" {
		if _, err := strconv.ParseUint(cfg.GuildID, 10, 64); err != nil {
			return nil, fmt.Errorf("GUILD_ID must be a numeric snowflake ID, got %q", cfg.GuildID)
		}
	}

	cfg.WarningInterval = parseDuration("WARNING_POLL_INTERVAL", 90*time.Second)
	cfg.TipsInterval = parseDuration("TIPS_POLL_INTERVAL", 5*time.Minute)
	cfg.StatusInterval = parseDuration("STATUS_POLL_INTERVAL", 10*time.Minute)

	return cfg, nil
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

func parseDuration(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}

func parseLogLevel(s string) slog.Level {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(s)); err != nil {
		return slog.LevelInfo
	}
	return lvl
}
