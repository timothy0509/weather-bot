package config

import (
	"testing"
	"time"
)

func TestLoadRequired(t *testing.T) {
	t.Setenv("DISCORD_TOKEN", "")
	t.Setenv("DISCORD_APPLICATION_ID", "")
	_, err := Load()
	if err == nil {
		t.Error("expected error for missing token")
	}
}

func TestLoadDefaults(t *testing.T) {
	t.Setenv("DISCORD_TOKEN", "test-token")
	t.Setenv("DISCORD_APPLICATION_ID", "123456")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.DatabasePath != "./weather-bot.db" {
		t.Errorf("database path = %q, want ./weather-bot.db", cfg.DatabasePath)
	}
	if cfg.WarningInterval != 90*time.Second {
		t.Errorf("warning interval = %v, want 90s", cfg.WarningInterval)
	}
	if cfg.TipsInterval != 5*time.Minute {
		t.Errorf("tips interval = %v, want 5m", cfg.TipsInterval)
	}
}

func TestLoadInvalidGuildID(t *testing.T) {
	t.Setenv("DISCORD_TOKEN", "tok")
	t.Setenv("DISCORD_APPLICATION_ID", "123")
	t.Setenv("GUILD_ID", "not-a-number")
	_, err := Load()
	if err == nil {
		t.Error("expected error for invalid guild ID")
	}
}
