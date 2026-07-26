package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadRequired(t *testing.T) {
	os.Clearenv()
	_, err := Load()
	if err == nil {
		t.Error("expected error for missing token")
	}
}

func TestLoadDefaults(t *testing.T) {
	os.Clearenv()
	os.Setenv("DISCORD_TOKEN", "test-token")
	os.Setenv("DISCORD_APPLICATION_ID", "123456")

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

func TestMustParseSnowflakeID(t *testing.T) {
	id := MustParseSnowflakeID("123456789")
	if id != 123456789 {
		t.Errorf("id = %d, want 123456789", id)
	}
}
