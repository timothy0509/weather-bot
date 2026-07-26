package db

import (
	"database/sql"
	"fmt"
	"time"
)

// GuildSettings stores per-guild configuration.
type GuildSettings struct {
	GuildID          string
	AlertChannelID   sql.NullString
	Language         string
	TideStation      string
	BotStatusEnabled bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// GetGuildSettings retrieves settings for a guild.
func (db *DB) GetGuildSettings(guildID string) (*GuildSettings, error) {
	gs := &GuildSettings{}
	err := db.QueryRow(
		`SELECT guild_id, alert_channel_id, language, tide_station, bot_status_enabled, created_at, updated_at
		 FROM guild_settings WHERE guild_id = ?`,
		guildID,
	).Scan(
		&gs.GuildID,
		&gs.AlertChannelID,
		&gs.Language,
		&gs.TideStation,
		&gs.BotStatusEnabled,
		&gs.CreatedAt,
		&gs.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		// Return default settings.
		return &GuildSettings{
			GuildID:          guildID,
			Language:         "en",
			TideStation:      "QUB",
			BotStatusEnabled: true,
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get guild settings: %w", err)
	}
	return gs, nil
}

// SetGuildSettings upserts settings for a guild.
func (db *DB) SetGuildSettings(gs *GuildSettings) error {
	_, err := db.Exec(
		`INSERT INTO guild_settings (guild_id, alert_channel_id, language, tide_station, bot_status_enabled, updated_at)
		 VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		 ON CONFLICT(guild_id) DO UPDATE SET
			 alert_channel_id = excluded.alert_channel_id,
			 language = excluded.language,
			 tide_station = excluded.tide_station,
			 bot_status_enabled = excluded.bot_status_enabled,
			 updated_at = CURRENT_TIMESTAMP`,
		gs.GuildID,
		gs.AlertChannelID,
		gs.Language,
		gs.TideStation,
		gs.BotStatusEnabled,
	)
	if err != nil {
		return fmt.Errorf("set guild settings: %w", err)
	}
	return nil
}

// SetLanguage updates the language for a guild.
func (db *DB) SetLanguage(guildID, language string) error {
	gs, err := db.GetGuildSettings(guildID)
	if err != nil {
		return err
	}
	gs.Language = language
	return db.SetGuildSettings(gs)
}

// SetAlertChannel updates the alert channel for a guild.
func (db *DB) SetAlertChannel(guildID, channelID string) error {
	gs, err := db.GetGuildSettings(guildID)
	if err != nil {
		return err
	}
	gs.AlertChannelID = sql.NullString{String: channelID, Valid: channelID != ""}
	return db.SetGuildSettings(gs)
}

// SetTideStation updates the default tide station for a guild.
func (db *DB) SetTideStation(guildID, station string) error {
	gs, err := db.GetGuildSettings(guildID)
	if err != nil {
		return err
	}
	gs.TideStation = station
	return db.SetGuildSettings(gs)
}

// SetBotStatusEnabled updates whether the bot status is enabled.
func (db *DB) SetBotStatusEnabled(guildID string, enabled bool) error {
	gs, err := db.GetGuildSettings(guildID)
	if err != nil {
		return err
	}
	gs.BotStatusEnabled = enabled
	return db.SetGuildSettings(gs)
}

// AllGuildSettings returns settings for all guilds.
func (db *DB) AllGuildSettings() ([]*GuildSettings, error) {
	rows, err := db.Query(`SELECT guild_id, alert_channel_id, language, tide_station, bot_status_enabled, created_at, updated_at FROM guild_settings`)
	if err != nil {
		return nil, fmt.Errorf("all guild settings: %w", err)
	}
	defer rows.Close()

	var results []*GuildSettings
	for rows.Next() {
		gs := &GuildSettings{}
		if err := rows.Scan(
			&gs.GuildID,
			&gs.AlertChannelID,
			&gs.Language,
			&gs.TideStation,
			&gs.BotStatusEnabled,
			&gs.CreatedAt,
			&gs.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan guild settings: %w", err)
		}
		results = append(results, gs)
	}
	return results, rows.Err()
}

// WarningState is a persisted warning record.
type WarningState struct {
	ID         int64
	Code       string
	Subtype    sql.NullString
	ActionCode string
	IssueTime  sql.NullString
	UpdateTime sql.NullString
	LastSeen   time.Time
}

// GetWarningState retrieves the most recent warning state by code.
func (db *DB) GetWarningState(code string) (*WarningState, error) {
	ws := &WarningState{}
	err := db.QueryRow(
		`SELECT id, code, subtype, action_code, issue_time, update_time, last_seen
		 FROM warning_state WHERE code = ? ORDER BY last_seen DESC LIMIT 1`,
		code,
	).Scan(
		&ws.ID,
		&ws.Code,
		&ws.Subtype,
		&ws.ActionCode,
		&ws.IssueTime,
		&ws.UpdateTime,
		&ws.LastSeen,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get warning state: %w", err)
	}
	return ws, nil
}

// LatestWarningStates returns the latest state for each warning code.
func (db *DB) LatestWarningStates() ([]*WarningState, error) {
	rows, err := db.Query(
		`SELECT id, code, subtype, action_code, issue_time, update_time, last_seen
		 FROM warning_state
		 WHERE id IN (
			 SELECT MAX(id) FROM warning_state GROUP BY code
		 )`,
	)
	if err != nil {
		return nil, fmt.Errorf("latest warning states: %w", err)
	}
	defer rows.Close()

	var results []*WarningState
	for rows.Next() {
		ws := &WarningState{}
		if err := rows.Scan(
			&ws.ID,
			&ws.Code,
			&ws.Subtype,
			&ws.ActionCode,
			&ws.IssueTime,
			&ws.UpdateTime,
			&ws.LastSeen,
		); err != nil {
			return nil, fmt.Errorf("scan warning state: %w", err)
		}
		results = append(results, ws)
	}
	return results, rows.Err()
}

// SaveWarningState upserts a warning state record.
func (db *DB) SaveWarningState(code, subtype, actionCode, issueTime, updateTime string) error {
	_, err := db.Exec(
		`INSERT INTO warning_state (code, subtype, action_code, issue_time, update_time, last_seen)
		 VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`,
		code, subtype, actionCode, issueTime, updateTime,
	)
	if err != nil {
		return fmt.Errorf("save warning state: %w", err)
	}
	return nil
}

// GetLatestTipsUpdateTime retrieves the latest tips update time.
func (db *DB) GetLatestTipsUpdateTime() (string, error) {
	var updateTime string
	err := db.QueryRow(
		`SELECT update_time FROM tips_state ORDER BY last_seen DESC LIMIT 1`,
	).Scan(&updateTime)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("get latest tips update time: %w", err)
	}
	return updateTime, nil
}

// SaveTipsUpdateTime records a tips update time.
func (db *DB) SaveTipsUpdateTime(updateTime string) error {
	_, err := db.Exec(
		`INSERT INTO tips_state (update_time, last_seen) VALUES (?, CURRENT_TIMESTAMP)`,
		updateTime,
	)
	if err != nil {
		return fmt.Errorf("save tips update time: %w", err)
	}
	return nil
}
