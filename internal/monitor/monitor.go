// Package monitor polls HKO for warnings and tips and sends alerts to Discord.
package monitor

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/time/rate"

	"weather-bot/internal/db"
	"weather-bot/internal/format"
	"weather-bot/internal/hko"
	"weather-bot/internal/i18n"
)

// Monitor runs background polling loops.
type Monitor struct {
	session *discordgo.Session
	db      *db.DB
	hko     *hko.Client
	logger  *slog.Logger
	cfg     Intervals
	wg      sync.WaitGroup
	limiter *rate.Limiter

	mu              sync.Mutex
	lastWarningPoll time.Time
	lastTipsPoll    time.Time
	lastStatusPoll  time.Time
}

// Intervals holds polling intervals.
type Intervals struct {
	Warning time.Duration
	Tips    time.Duration
	Status  time.Duration
}

// New creates a new Monitor.
func New(session *discordgo.Session, database *db.DB, client *hko.Client, logger *slog.Logger, cfg Intervals) *Monitor {
	return &Monitor{
		session: session,
		db:      database,
		hko:     client,
		logger:  logger,
		cfg:     cfg,
		limiter: rate.NewLimiter(1, 3),
	}
}

// Start begins all background loops.
func (m *Monitor) Start(ctx context.Context) {
	m.wg.Add(3)
	go func() {
		defer m.wg.Done()
		m.warningLoop(ctx)
	}()
	go func() {
		defer m.wg.Done()
		m.tipsLoop(ctx)
	}()
	go func() {
		defer m.wg.Done()
		m.statusLoop(ctx)
	}()
}

// Wait blocks until all background loops have stopped.
func (m *Monitor) Wait() {
	m.wg.Wait()
}

func (m *Monitor) warningLoop(ctx context.Context) {
	ticker := time.NewTicker(m.cfg.Warning)
	defer ticker.Stop()
	cleanupTicker := time.NewTicker(24 * time.Hour)
	defer cleanupTicker.Stop()
	m.checkWarnings()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.checkWarnings()
		case <-cleanupTicker.C:
			if err := m.db.CleanupOldStates(); err != nil {
				m.logger.Warn("state cleanup failed", slog.Any("err", err))
			}
		}
	}
}

func (m *Monitor) tipsLoop(ctx context.Context) {
	ticker := time.NewTicker(m.cfg.Tips)
	defer ticker.Stop()
	m.checkTips()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.checkTips()
		}
	}
}

func (m *Monitor) statusLoop(ctx context.Context) {
	ticker := time.NewTicker(m.cfg.Status)
	defer ticker.Stop()
	m.updateStatus()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.updateStatus()
		}
	}
}

// Health returns the last successful poll times for each loop.
func (m *Monitor) Health() (warning, tips, status time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.lastWarningPoll, m.lastTipsPoll, m.lastStatusPoll
}

func (m *Monitor) checkWarnings() {
	ws, err := m.hko.GetWarningSummary("en")
	if err != nil {
		m.logger.Warn("warning summary poll failed", slog.Any("err", err))
		return
	}
	m.mu.Lock()
	m.lastWarningPoll = time.Now()
	m.mu.Unlock()

	for code, w := range ws {
		state, err := m.db.GetWarningState(code)
		if err != nil {
			m.logger.Warn("failed to get warning state", slog.Any("err", err))
			continue
		}

		actionCode := w.ActionCode
		if actionCode == "" {
			actionCode = "Active"
		}

		shouldAlert := false
		alertType := "warning_issued"
		if state == nil {
			shouldAlert = true
		} else if state.ActionCode == "CANCEL" && actionCode != "CANCEL" {
			shouldAlert = true
			alertType = "warning_issued"
		} else if state.ActionCode != actionCode || state.Subtype.String != w.Subtype {
			shouldAlert = true
			if actionCode == "CANCEL" {
				alertType = "warning_cancelled"
			} else {
				alertType = "warning_updated"
			}
		}

		if shouldAlert {
			m.sendWarningAlert(code, w, alertType)
		}

		if err := m.db.SaveWarningState(code, w.Subtype, actionCode, w.IssueTime, w.UpdateTime); err != nil {
			m.logger.Warn("failed to save warning state", slog.Any("err", err))
		}
	}

	// Detect cancellations: active codes in DB but not in current summary.
	latestStates, err := m.db.LatestWarningStates()
	if err != nil {
		m.logger.Warn("failed to get latest warning states", slog.Any("err", err))
		return
	}
	for _, state := range latestStates {
		if state.ActionCode == "CANCEL" {
			continue
		}
		if _, active := ws[state.Code]; !active {
			m.sendWarningAlert(state.Code, &hko.WarningSummaryWarning{
				Name:       state.Code,
				Code:       state.Code,
				Type:       "Cancelled",
				ActionCode: "CANCEL",
				UpdateTime: state.UpdateTime.String,
			}, "warning_cancelled")
			if err := m.db.SaveWarningState(state.Code, state.Subtype.String, "CANCEL", state.IssueTime.String, state.UpdateTime.String); err != nil {
				m.logger.Warn("failed to save cancellation state", slog.Any("err", err))
			}
		}
	}
}

func (m *Monitor) sendWarningAlert(code string, w *hko.WarningSummaryWarning, alertType string) {
	settings, err := m.db.AllGuildSettings()
	if err != nil {
		m.logger.Warn("failed to get guild settings for alert", slog.Any("err", err))
		return
	}

	for _, gs := range settings {
		if !gs.AlertChannelID.Valid {
			continue
		}
		lang := i18n.Normalize(gs.Language)

		title := i18n.T(alertType, lang)
		color := 0xE74C3C
		if alertType == "warning_cancelled" {
			color = 0x95A5A6
		}

		description := w.Type
		info, err := m.hko.GetWarningInfo(string(lang))
		if err == nil {
			var contents []string
			for _, d := range info.Details {
				if strings.HasPrefix(d.WarningStatementCode, code) {
					contents = append(contents, d.Contents...)
				}
			}
			if len(contents) > 0 {
				description = strings.Join(contents, "\n\n")
			}
		}

		footerText := ""
		if alertType == "warning_cancelled" {
			footerText = fmt.Sprintf("%s %s", i18n.T("cancelled_at", lang), format.FormatTime(w.UpdateTime))
		} else {
			footerText = fmt.Sprintf("%s %s", i18n.T("in_effect_since", lang), format.FormatTime(w.IssueTime))
		}

		embed := &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("%s: %s", title, w.Name),
			Description: description,
			Color:       color,
			Footer: &discordgo.MessageEmbedFooter{
				Text: footerText,
			},
		}

		if err := m.limiter.Wait(context.Background()); err != nil {
			return
		}
		_, err = m.session.ChannelMessageSendEmbed(gs.AlertChannelID.String, embed)
		if err != nil {
			m.logger.Warn("failed to send warning alert", slog.Any("err", err))
		}
	}
}

func (m *Monitor) checkTips() {
	tips, err := m.hko.GetSpecialWeatherTips("en")
	if err != nil {
		m.logger.Warn("tips poll failed", slog.Any("err", err))
		return
	}
	m.mu.Lock()
	m.lastTipsPoll = time.Now()
	m.mu.Unlock()
	if len(tips.SWT) == 0 {
		return
	}

	latest := tips.SWT[0]
	lastUpdate, err := m.db.GetLatestTipsUpdateTime()
	if err != nil {
		m.logger.Warn("failed to get latest tips update time", slog.Any("err", err))
		return
	}
	if latest.UpdateTime == lastUpdate {
		return
	}

	settings, err := m.db.AllGuildSettings()
	if err != nil {
		m.logger.Warn("failed to get guild settings for tips", slog.Any("err", err))
		return
	}
	for _, gs := range settings {
		if !gs.AlertChannelID.Valid {
			continue
		}
		lang := i18n.Normalize(gs.Language)

		tipsLocalized, err := m.hko.GetSpecialWeatherTips(string(lang))
		desc := latest.Desc
		if err == nil && len(tipsLocalized.SWT) > 0 {
			desc = tipsLocalized.SWT[0].Desc
		}

		embed := &discordgo.MessageEmbed{
			Title:       i18n.T("special_tips", lang),
			Description: desc,
			Color:       0xF39C12,
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("%s: %s", i18n.T("updated_at", lang), format.FormatTime(latest.UpdateTime)),
			},
		}

		if err := m.limiter.Wait(context.Background()); err != nil {
			return
		}
		_, err = m.session.ChannelMessageSendEmbed(gs.AlertChannelID.String, embed)
		if err != nil {
			m.logger.Warn("failed to send tips alert", slog.Any("err", err))
		}
	}

	if err := m.db.SaveTipsUpdateTime(latest.UpdateTime); err != nil {
		m.logger.Warn("failed to save tips update time", slog.Any("err", err))
	}
}

func (m *Monitor) updateStatus() {
	// Check only first guild with status enabled. For global bot status,
	// just use the current weather once any guild has it enabled.
	settings, err := m.db.AllGuildSettings()
	if err != nil {
		m.logger.Warn("failed to get guild settings for status", slog.Any("err", err))
		return
	}

	statusEnabled := false
	for _, gs := range settings {
		if gs.BotStatusEnabled {
			statusEnabled = true
			break
		}
	}
	if !statusEnabled {
		return
	}

	w, err := m.hko.GetCurrentWeather("en")
	if err != nil {
		m.logger.Warn("status update weather fetch failed", slog.Any("err", err))
		return
	}
	m.mu.Lock()
	m.lastStatusPoll = time.Now()
	m.mu.Unlock()

	var status string
	if r, ok := hko.ReadingByPlace(w.Temperature, "Hong Kong Observatory"); ok {
		status = fmt.Sprintf("%s %.1f°C", r.Place, r.Value)
	} else if len(w.Temperature.Data) > 0 {
		status = fmt.Sprintf("%.1f°C", w.Temperature.Data[0].Value)
	} else {
		status = "HK Weather"
	}

	if len(status) > 128 {
		status = status[:128]
	}

	if err := m.session.UpdateGameStatus(0, status); err != nil {
		m.logger.Warn("failed to update bot status", slog.Any("err", err))
	}
}
