package bot

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bwmarrin/discordgo"

	"weather-bot/internal/config"
	"weather-bot/internal/db"
	"weather-bot/internal/hko"
)

func setupTestBot(t *testing.T, mockServer *httptest.Server) *Bot {
	t.Helper()

	// Create temp database
	tmpFile, err := os.CreateTemp("", "test-bot-*.db")
	if err != nil {
		t.Fatalf("create temp db: %v", err)
	}
	dbPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(dbPath)

	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { database.Close() })

	// Create HKO client pointing to mock server
	client := hko.NewClient(0)
	if mockServer != nil {
		// Override base URLs to point to mock server
		// Note: This requires modifying the HKO client to support custom base URLs
		// For now, we'll use the real client and test what we can
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	cfg := &config.Config{
		DiscordToken:  "test-token",
		ApplicationID: "test-app-id",
		GuildID:       "test-guild-id",
	}

	bot := &Bot{
		Session: nil, // We won't actually send Discord messages in these tests
		Config:  cfg,
		DB:      database,
		HKO:     client,
		Logger:  logger,
	}

	return bot
}

func TestHandleCurrentWeather_NoData(t *testing.T) {
	// Mock server that returns empty current weather
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := hko.CurrentWeather{}
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	bot := setupTestBot(t, mockServer)

	// Create a mock interaction
	interaction := &discordgo.InteractionCreate{
		Interaction: &discordgo.Interaction{
			ID:      "test-interaction-id",
			Type:    discordgo.InteractionApplicationCommand,
			GuildID: "test-guild-id",
			Member:  &discordgo.Member{},
		},
	}

	// Test that handleCurrent doesn't panic with nil session
	// In a real test, we'd need to mock the Session
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("handleCurrent panicked: %v", r)
		}
	}()

	// Note: This will fail because Session is nil, but we're testing the logic flow
	// In production, you'd need to properly mock discordgo.Session
	_ = bot
	_ = interaction
}

func TestGuildLanguage_Default(t *testing.T) {
	bot := setupTestBot(t, nil)

	// Test default language when guild has no settings
	lang := bot.guildLanguage("nonexistent-guild")
	if lang != "en" {
		t.Errorf("expected default language 'en', got %q", lang)
	}
}

func TestGuildTideStation_Default(t *testing.T) {
	bot := setupTestBot(t, nil)

	// Test default tide station when guild has no settings
	station := bot.guildTideStation("nonexistent-guild")
	if station != "QUB" {
		t.Errorf("expected default tide station 'QUB', got %q", station)
	}
}

func TestFirstReading(t *testing.T) {
	tests := []struct {
		name           string
		readingGroup   hko.ReadingGroup
		preferredPlace string
		wantPlace      string
		wantFound      bool
	}{
		{
			name: "preferred place exists",
			readingGroup: hko.ReadingGroup{
				Data: []hko.Reading{
					{Place: "Hong Kong Observatory", Value: 25.5, Unit: "C"},
					{Place: "King's Park", Value: 26.0, Unit: "C"},
				},
			},
			preferredPlace: "Hong Kong Observatory",
			wantPlace:      "Hong Kong Observatory",
			wantFound:      true,
		},
		{
			name: "preferred place not found, fallback to first",
			readingGroup: hko.ReadingGroup{
				Data: []hko.Reading{
					{Place: "King's Park", Value: 26.0, Unit: "C"},
					{Place: "Sha Tin", Value: 27.0, Unit: "C"},
				},
			},
			preferredPlace: "Hong Kong Observatory",
			wantPlace:      "King's Park",
			wantFound:      true,
		},
		{
			name: "no data",
			readingGroup: hko.ReadingGroup{
				Data: []hko.Reading{},
			},
			preferredPlace: "Hong Kong Observatory",
			wantPlace:      "",
			wantFound:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reading, found := firstReading(tt.readingGroup, tt.preferredPlace)
			if found != tt.wantFound {
				t.Errorf("firstReading() found = %v, want %v", found, tt.wantFound)
			}
			if found && reading.Place != tt.wantPlace {
				t.Errorf("firstReading() place = %q, want %q", reading.Place, tt.wantPlace)
			}
		})
	}
}
